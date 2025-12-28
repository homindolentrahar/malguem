package command

import (
	"fmt"
	"malguem/internal/model"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var new = &cobra.Command{
	Use:   "new",
	Short: "Create new template",
	Long:  "Create new template for your boilerplate code",
	Run: func(cmd *cobra.Command, args []string) {
		var name string

		// Check if there's any argument for "new" command
		if len(args) < 1 {
			fmt.Println("Please provide name after command")
			os.Exit(1)
		}

		name = args[0]

		// Create folder based on accepted `name`
		templatePath := filepath.Join("templates", name)
		os.MkdirAll(templatePath, os.ModePerm)

		// Create template.yaml file
		err := createContract(name)
		if err != nil {
			fmt.Printf("Failed to create a new template: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🚀  %s template created. Happy coding !", name)
	},
}

func createContract(template string) error {
	// Create contract.yaml file inside template's folder
	filename := "contract.yaml"
	filePath := fmt.Sprintf("templates/%s/%s", template, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Define the initial contract file
	contract := model.Contract{
		Name: template,
		Variables: map[string]model.Variable{
			"name": {
				Type:    "string",
				Prompt:  "Enter your name",
				Default: template,
			},
		},
	}

	// Populate contract into YAML file
	data, err := yaml.Marshal(&contract)
	if err != nil {
		return err
	}

	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
