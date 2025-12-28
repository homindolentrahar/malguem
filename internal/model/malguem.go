package model

type Malguem struct {
	Name      string              `yaml:"name"`
	Templates map[string]Template `yaml:"templates"`
}

type Template struct {
	Path   string  `yaml:"path,omitempty"`
	Remote *Remote `yaml:"remote,omitempty"`
	Output string  `yaml:"output,omitempty"`
}

type Remote struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
	Ref  string `yaml:"ref"`
}
