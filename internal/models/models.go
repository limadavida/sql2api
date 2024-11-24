package models

type ProjectConfig struct {
	Project  string
	RootDir  string
	Version  string
	Author   string
	Type     string
	Servers  string
	Database DatabaseConfig
	Models   map[string]ModelFiles
}

type DatabaseConfig struct {
	Type string
	Name string
}

type ModelFiles struct {
	File []string
}
