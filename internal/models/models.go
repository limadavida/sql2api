package models

type ProjectConfig struct {
	Project     string
	RootDir     string
	Version     string
	Author      string
	Type        string
	Servers     string
	Database    DatabaseConfig
	Models      map[string]string
	ModelsSetup map[string]ModelConfig
}

type DatabaseConfig struct {
	Type string
	Name string
}

type ModelFiles struct {
	File []string
}

type ModelConfig struct {
	Description string
	HTTPMethod  string
	Async       bool
	Parameters  string
	Query       string
}
