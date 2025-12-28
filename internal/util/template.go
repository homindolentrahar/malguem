package util

import (
	"fmt"
	"io/fs"
	"malguem/internal/config"
	"malguem/internal/prompt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/cbroglie/mustache"
)

var sectionMustachePattern = regexp.MustCompile(`{{#(\w+)}}(\w+){{/(\w+)}}`)
var pathMustachePattern = regexp.MustCompile(`{{\s*(\w+)\s*}}`)

func RenderTemplate(source, output string) error {
	// Trim prefixes of './' inside the
	source = strings.TrimPrefix(source, "./")

	// Read config `contract.yaml`
	contractPath := filepath.Join(source, "contract.yaml")
	contract, err := config.ReadContract(contractPath)
	if err != nil {
		return err
	}

	// Read variables from contract
	var inputs = make(map[string]string)
	for key, value := range contract.Variables {
		// Prompt user input for value
		inputPrompt := prompt.PromptInput(fmt.Sprintf("%s (%s)", value.Prompt, value.Default))
		if inputPrompt == "" {
			inputPrompt = value.Default.(string)
		}

		inputs[key] = inputPrompt
	}

	// Make sure the `output` directory exists
	err = os.MkdirAll(output, os.ModePerm)
	if err != nil {
		return err
	}

	// Walk every file and folder in `source` directory
	return filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip to render root directory
		// Skip to render `contract.yaml` file
		if path == source || filepath.Base(path) == "contract.yaml" {
			return nil
		}

		processedPath := ""

		// Render mustache template in the path
		if info.IsDir() {
			processedPath = preprocessPath(path, inputs, contract.DirCase)
		} else {
			// Split path into directory and file segments
			dirSegment := filepath.Dir(path)
			fileSegment := filepath.Base(path)

			// Preprocess directory path and file path first
			dirPath := preprocessPath(dirSegment, inputs, contract.DirCase)
			filePath := preprocessPath(fileSegment, inputs, contract.FileCase)

			// Join `dirPath` and `filePath` together to form complete path
			fullpath := filepath.Join(dirPath, filePath)
			processedPath = fullpath
		}

		// Trim prefixes from `path` in `processedPath`, so it only yields variable name
		processedPath = strings.TrimPrefix(processedPath, source)

		// Check if walked info is directory
		// Then make sure to create directory
		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(output, processedPath), os.ModePerm)
			if err != nil {
				return err
			}

			return nil
		}

		// Render mustache template of the content
		outputPath := filepath.Join(output, processedPath)
		err = renderMustacheTemplate(path, outputPath, inputs)
		if err != nil {
			return err
		}

		return nil
	})
}

func renderMustacheTemplate(source, output string, data map[string]string) error {
	// Read content from source path
	file, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	// Process to get variable from mustache template in content
	fileString := string(file)
	processedTemplate := preprocessContent(fileString, data)

	// Render processed mustache template
	result, err := mustache.Render(processedTemplate, data)
	if err != nil {
		return err
	}

	return os.WriteFile(output, []byte(result), os.FileMode(0644))
}

func preprocessPath(path string, data map[string]string, typeCase string) string {
	return pathMustachePattern.ReplaceAllStringFunc(path, func(match string) string {
		variable := pathMustachePattern.FindStringSubmatch(match)[1]

		// If the `typeCase` is empty, then fallback to "snake_case"
		if typeCase == "" {
			typeCase = Snake
		}

		value, isExists := data[variable]
		if !isExists {
			return ""
		}

		return formatCase(value, typeCase)
	})
}

func preprocessContent(contentString string, data map[string]string) string {
	return sectionMustachePattern.ReplaceAllStringFunc(contentString, func(match string) string {
		matches := sectionMustachePattern.FindStringSubmatch(match)
		// {{#format}} varName {{/format}}
		//     1         2			3
		// NOTE: 0 will give the whole unformatted match

		// Invalid format, since we only capture 3 group of submatch
		if len(matches) > 4 {
			return match
		}

		// Extract format and variable from mustache template
		firstTagFormat, variable, secondTagFormat := matches[1], matches[2], matches[3]

		if !slices.Contains(TypeCases, firstTagFormat) || firstTagFormat != secondTagFormat {
			fmt.Printf("\nFormat tag invalid, please check your template again. For correct template tag, please visit ...\n")
			// Rollout the generated code before exit
			os.Exit(1)
		}

		value, isExists := data[variable]
		if !isExists {
			return match
		}

		return formatCase(value, firstTagFormat)
	})
}

func formatCase(value, typeCase string) string {
	switch typeCase {
	case Pascal:
		return PascalCase(value)
	case Camel:
		return CamelCase(value)
	case Snake:
		return SnakeCase(value)
	case Kebab:
		return KebabCase(value)
	case Title:
		return TitleCase(value)
	default:
		return value
	}
}
