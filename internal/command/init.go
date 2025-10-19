package command

import (
	"fmt"
	"malguem/internal/model"
	"malguem/internal/util"
	"os"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var initialize = &cobra.Command{
	Use:   "init",
	Short: "Initialize malguem in the project",
	Run: func(cmd *cobra.Command, args []string) {
		createMalguemFile()
		fmt.Printf("🌤️  My name is Malguem, let's generate some code!")
	},
}

func createMalguemFile() error {
	// Create malguem.yaml file
	file, err := os.Create("malguem.yaml")
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
