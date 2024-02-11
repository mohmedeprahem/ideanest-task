package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type appConfigYaml struct {
	// Define your struct fields to match the structure of your YAML data
	Jwt struct {
		AtSecret string `yaml:"atSecret"`
		RtSecret string `yaml:"rtSecret"`
	}
	Redis struct {
		Address string `yaml:"address"`
		Password string `yaml:"password"`
		Db int `yaml:"db"`
	}
}

func ReadAppConfig() (*appConfigYaml, error) {
	var config appConfigYaml

	file, err := os.ReadFile("config/app-config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

