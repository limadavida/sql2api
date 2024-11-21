package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

var ConfigData *Config

func init() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	ConfigData = config

	yamlConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Erro ao converter configuração para YAML: %v", err)
	}

	log.Printf("Configuração carregada: \n%s", string(yamlConfig))
}
