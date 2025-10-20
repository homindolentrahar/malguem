package model

type Contract struct {
	Name      string              `yaml:"name"`
	FileCase  string              `yaml:"file_case,omitempty"`
	DirCase   string              `yaml:"dir_case,omitempty"`
	Variables map[string]Variable `yaml:"variables"`
}

type Variable struct {
	Type    any    `yaml:"type"`
	Prompt  string `yaml:"prompt"`
	Default any    `yaml:"default"`
}
