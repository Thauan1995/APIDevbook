package config

import (
	"context"
	"fmt"
	"site/utils/consts"
	"site/utils/log"
	"strings"

	"cloud.google.com/go/datastore"
)

const (
	ConfigKind = "Config"

	APIToken = "apitoken"

	ChaveAutenticacaoAcesso = "login.chaveautenticacao"
	TempoExpiracaoAcesso    = "login.tempoexpiracao"

	ElasticSearchEndpoint = "elasticsearch.endpoint"
	ElasticSearchUsername = "elasticsearch.username"
	ElasticSearchPassword = "elasticsearch.password"
)

var SecretKey []byte

type Config struct {
	Name  string `datastore:"-"`
	Value string `datastore:",noindex"`
}

func GetDefault(c context.Context, configName string, valDefault string) *Config {
	config, err := GetConfig(c, configName)
	if err != nil || config == nil {
		config = &Config{Name: configName, Value: valDefault}
	}

	return config
}

func GetConfig(c context.Context, configName string) (*Config, error) {

	key := datastore.NameKey(ConfigKind, configName, nil)

	var config Config
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return &config, err
	}
	defer datastoreClient.Close()
	err = datastoreClient.Get(c, key, &config)
	if err == nil {
		config.Name = configName
	}

	return &config, err
}

func ListConfigs(c context.Context) ([]Config, error) {
	var confs []Config

	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return confs, err
	}
	defer datastoreClient.Close()
	q := datastore.NewQuery(ConfigKind)
	keys, err := datastoreClient.GetAll(c, q, &confs)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}

		if !strings.Contains(err.Error(), "no such struct field") {
			return nil, err
		}
	}

	for index, key := range keys {
		confs[index].Name = key.Name
	}

	return confs, nil
}

func PutConfig(c context.Context, config *Config) error {

	if config == nil {
		return fmt.Errorf("A nil config cannot be saved")
	}

	if config.Name == "" {
		return fmt.Errorf("The 'config.name' cannot be empty")
	}

	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()
	key := datastore.NameKey(ConfigKind, config.Name, nil)
	key, err = datastoreClient.Put(c, key, config)
	if err != nil {
		return err
	}

	config.Name = key.Name
	return nil
}

func BuscaSecret(c context.Context) {
	chave, err := GetConfig(c, ChaveAutenticacaoAcesso)
	if err != nil {
		log.Warningf(c, "Falha ao buscar chave de autenticação %v", err)
		return
	}

	SecretKey = []byte(chave.Value)
}
