package config

import "github.com/spf13/viper"

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
