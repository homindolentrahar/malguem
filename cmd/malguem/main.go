package main

import (
	"fmt"
	"malguem/internal/config"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	appConfig, err := config.Read()
	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("Name: %s\nDesc: %s\nVersion: %s\n", appConfig.Name, appConfig.Desc, appConfig.Version))
}
