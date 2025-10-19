package model

type Malguem struct {
	Name      string              `yaml:"name"`
	Templates map[string]Template `yaml:"templates"`
}

type Template struct {
	Path   string  `yaml:"path,omitempty"`
	Github *Source `yaml:"github,omitempty"`
	Output string  `yaml:"output,omitempty"`
}

type Source struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
	Ref  string `yaml:"ref"`
}
