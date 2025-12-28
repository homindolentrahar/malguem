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

func Read() (*Info, error) {
	var info Info

	err := viper.Unmarshal(&info)
	if err != nil {
		return nil, err
	}

	return &info, err
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

func ReadContract(path string) (*model.Contract, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var contract model.Contract
	err = yaml.Unmarshal(file, &contract)
	if err != nil {
		return nil, err
	}

	return &contract, nil
}
