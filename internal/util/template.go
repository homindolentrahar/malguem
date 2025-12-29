package util

import (
	"fmt"
	"io/fs"
	"malguem/internal/config"
	"malguem/internal/model"
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
	// todo: add these into separate validateConfig in the make.go
	source = strings.TrimPrefix(source, "./")
	output = strings.TrimPrefix(output, "./")

	// Validate contract
	contract, err := validateContract(source)
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

	// Create the temp dir to holds rendered result
	tempDir, err := os.MkdirTemp(".", "malguem-make-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Walk every file and folder in `source` directory
	err = filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip to render root directory
		// also skip to render `contract.yaml` file
		if path == source || filepath.Base(path) == "contract.yaml" {
			return nil
		}

		// Validate template content first
		err = validateTemplate(source, path, inputs, info.IsDir())
		if err != nil {
			return err
		}

		processedPath := ""

		// Format file path with proper type cases before rendering the template
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

		//Check if walked info is directory
		//Then make sure to create directory
		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(tempDir, processedPath), os.ModePerm)
			if err != nil {
				return err
			}

			return nil
		}

		// Render mustache template of the content into `tempDir`
		outputPath := filepath.Join(tempDir, processedPath)
		err = renderMustacheTemplate(path, outputPath, inputs)
		if err != nil {
			return err
		}

		// Validate rendered result
		err = validateResult(outputPath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Commit the template rendering
	return commit(tempDir, output)
}

// Validate contract file of `contract.yaml`
func validateContract(source string) (*model.Contract, error) {
	// Ensure `contract.yaml` is present in `source` path
	contractPath := filepath.Join(source, "contract.yaml")
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("`contract.yaml` file does not exist in '%s'", source)
	}
	contract, err := config.ReadContract(contractPath)
	if err != nil {
		return nil, err
	}

	// Ensure `contract.yaml` has name
	if contract.Name == "" {
		return nil, fmt.Errorf("`name` is empty in the contract")
	}

	// Ensure `contract.yaml` has proper `file_case`
	if !slices.Contains(TypeCases, contract.FileCase) {
		return nil, fmt.Errorf("`file_case` in contract is invalid")
	}

	// Ensure `contract.yaml` has proper `dir_case`
	if !slices.Contains(TypeCases, contract.DirCase) {
		return nil, fmt.Errorf("`dir_case` in contract is invalid")
	}

	// Ensure `contract.yaml` has empty `variables`
	if len(contract.Variables) < 1 {
		return nil, fmt.Errorf("contract does not have any `variables`")
	}

	return contract, nil
}

// Validate mustache template in content
func validateTemplate(base, path string, data map[string]string, isDirectory bool) error {
	// Get relative path of the file
	relativePath, err := filepath.Rel(base, path)
	if err != nil {
		return err
	}
	// Return error if there's unsafe path
	if strings.HasPrefix(relativePath, "..") || strings.HasPrefix(relativePath, ".") {
		return fmt.Errorf("unsafe path in %s", relativePath)
	}
	// Check if `relativePath` has mustache opening and closing tag
	if strings.Contains(relativePath, "{{") && strings.Contains(relativePath, "}}") {
		// Validate mustache syntax
		return ValidateMustacheSyntax(path, data, isDirectory)
	}

	return nil
}

// Validate rendered result
func validateResult(output string) error {
	//No path traversal
	fileName := filepath.Base(output)
	if strings.Contains(fileName, "{{") || strings.Contains(fileName, "}}") {
		return fmt.Errorf("unrendered mustache template in path, check your template again\n")
	}
	// Check if there's unrendered {{ or }} inside file's content
	file, err := os.ReadFile(output)
	if err != nil {
		return err
	}
	if mustacheBlockPattern.Match(file) {
		return fmt.Errorf("unrendered mustache template inside content, check your template again\n")
	}

	return nil
}

// Commit rendered template in `temp` into `output`
func commit(temp, output string) error {
	parentOutput := filepath.Dir(output)
	// Make sure the `parentOutput` directory exists first
	os.MkdirAll(parentOutput, os.FileMode(0755))

	// Check if `output` is already exists, if it does then remove the old one before renaming
	if _, err := os.Stat(output); err == nil {
		os.RemoveAll(output)
	}
	return os.Rename(temp, output)
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
		tag, variable := matches[1], matches[2]

		value, isExists := data[variable]
		if !isExists {
			return match
		}

		return formatCase(value, tag)
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
