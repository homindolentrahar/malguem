package command

import (
	"fmt"
	"malguem/internal/model"
	"malguem/internal/util"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var initialize = &cobra.Command{
	Use:   "init",
	Short: "Initialize malguem in the project",
	Run: func(cmd *cobra.Command, args []string) {
		var projectName *string
		if len(args) > 0 {
			projectName = &args[0]
		}

		err := createMalguemFile(projectName)
		if err != nil {
			fmt.Printf("Failed to initialize malguem inside your project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🌤️  I'm malguem, let's generate some code!")
	},
}

func createMalguemFile(projectName *string) error {
	// Create malguem.yaml file
	malguemPath := "malguem.yaml"
	if projectName != nil {
		// Create new directory
		os.MkdirAll(*projectName, os.ModePerm)
		malguemPath = filepath.Join(*projectName, malguemPath)
	}

	file, err := os.Create(malguemPath)
	if err != nil {
		return err
	}
	defer file.Close()

	currentDir, err := util.CurrentDir()
	if err != nil {
		return err
	}

	// Define the initial malguem config
	malguem := model.Malguem{
		Name: currentDir,
		Templates: map[string]model.Template{
			"example": {
				Path:   "./example/path",
				Output: "./template/output",
			},
		},
	}

	// Populate config into YAML file
	data, err := yaml.Marshal(&malguem)
	if err != nil {
		return err
	}

	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
