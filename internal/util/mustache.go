package util

import (
	"fmt"
	"os"
	"regexp"
	"slices"
)

var mustacheBlockPattern = regexp.MustCompile(`{{\s*([#^/])\s*([a-zA-Z0-9_.-]+)\s*}}`)

type MustacheTokenType int

const (
	Open MustacheTokenType = iota
	Close
)

type MustacheToken struct {
	Type     MustacheTokenType
	Tag      string
	Sigil    string
	Position int
}

// ValidateMustacheSyntax Validate mustache syntax
func ValidateMustacheSyntax(path string, data map[string]string, isDirectory bool) error {
	var content string

	if isDirectory {
		content = path
	} else {
		// Open content of the file
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content = string(file)
	}

	// Slice for storing open token type
	var stacks []MustacheToken
	tokens := tokenizeMustache(content)

	for _, token := range tokens {
		// If token type is `Open`, then add token into `stacks`
		if token.Type == Open {
			stacks = append(stacks, token)
			continue
		}

		// Check if `stacks` is empty, means that tag is closed without any opening tag
		if len(stacks) == 0 {
			return fmt.Errorf("tag %q closed without any opening tag in %s at position %d\n", token.Tag, path, token.Position)
		}

		// Pop the last opening tag inside `stacks`, means that tag is closed properly
		lastToken := stacks[len(stacks)-1]
		stacks = stacks[:len(stacks)-1]

		// Check if the last tag is not equal with current tag, means that the order of tag is not valid
		if lastToken.Tag != token.Tag {
			return fmt.Errorf(
				"tag mismatch: opened with %q but closed with %q in %s at position %d\n",
				lastToken.Tag, token.Tag, path, token.Position,
			)
		}

		// Check if the tag is not a valid type cases nor defined variable in `data`, means that tag is unknown
		_, isValidData := data[token.Tag]
		if !slices.Contains(TypeCases, token.Tag) && !isValidData {
			return fmt.Errorf(
				"unknown tag %q, not a valid type cases nor defined in `contract.yaml`\n",
				token.Tag,
			)
		}
	}

	// Check if stacks is not empty, means that tag is not closed
	if len(stacks) > 0 {
		lastToken := stacks[len(stacks)-1]
		return fmt.Errorf("unclosed tag %q in mustache block in %s at position %d \n", lastToken.Tag, path, lastToken.Position)
	}

	return nil
}

func tokenizeMustache(content string) []MustacheToken {
	matches := mustacheBlockPattern.FindAllStringSubmatchIndex(content, -1)
	tokens := make([]MustacheToken, 0, len(matches))

	for _, match := range matches {
		sigil := content[match[2]:match[3]]
		tag := content[match[4]:match[5]]

		tokenType := Open
		if sigil == "/" {
			tokenType = Close
		}

		tokens = append(tokens, MustacheToken{
			Type:     tokenType,
			Tag:      tag,
			Sigil:    sigil,
			Position: match[0],
		})
	}

	return tokens
}
