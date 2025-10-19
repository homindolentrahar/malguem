package main

import (
	"malguem/internal/command"
	"malguem/internal/config"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	_, err = config.Read()
	if err != nil {
		panic(err)
	}

	command.Run()
}
