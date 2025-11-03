package config

import (
	"malguem/internal/model"
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func Init() error {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}

func Read() (*Model, error) {
	var model Model

	err := viper.Unmarshal(&model)
	if err != nil {
		return nil, err
	}

	return &model, err
}

func ReadMalguem() (*model.Malguem, error) {
	file, err := os.ReadFile("malguem.yaml")
	if err != nil {
		return nil, err
	}

	var malguem model.Malguem
	err = yaml.Unmarshal(file, &malguem)
	if err != nil {
		return nil, err
	}

	return &malguem, nil
}
