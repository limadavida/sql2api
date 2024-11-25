package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/limadavida/sql2api/internal/database"
	"github.com/limadavida/sql2api/internal/logger"
	"github.com/limadavida/sql2api/internal/models"
	"github.com/limadavida/sql2api/internal/utils"
)

var ()

type Config struct {
	ModelsPath string
	TablesPath string
}

func NewConfig(folder, projectName string) *Config {
	projectBasePath := filepath.Join(folder, projectName)
	return &Config{
		ModelsPath: filepath.Join(projectBasePath, "models"),
		TablesPath: filepath.Join(projectBasePath, "tables"),
	}
}

func (c *Config) getModels() (map[string]string, error) {
	content := make(map[string]string)
	modelsFiles, err := utils.ListFiles(c.ModelsPath, ".sql")
	if err != nil {
		return nil, err
	}

	for _, modelsFile := range modelsFiles {
		content[modelsFile] = utils.ReadFile(modelsFile)
	}
	return content, nil
}

func (c *Config) getTables() (map[string]string, error) {
	content := make(map[string]string)
	modelsFiles, err := utils.ListFiles(c.TablesPath, ".sql")
	if err != nil {
		return nil, err
	}

	for _, modelsFile := range modelsFiles {
		content[modelsFile] = utils.ReadFile(modelsFile)
	}
	return content, nil
}

func (c *Config) GetProject() (map[string]string, map[string]string, error) {
	var m map[string]string
	var t map[string]string
	var err error

	m, err = c.getModels()
	if err != nil {
		return nil, nil, err
	}

	t, err = c.getTables()
	if err != nil {
		return nil, nil, err
	}
	return m, t, nil
}

func (c *Config) CreateTables(db database.Database, tables map[string]string) error {
	logger := logger.GetLogger()

	for tableName, tableSql := range tables {
		err := db.Execute(tableSql)
		if err != nil {
			return err
		}
		logger.Info("Tabela", tableName, "criada com sucesso!")
	}
	return nil
}

func (c *Config) CreateProjectPath() error {
	logger := logger.GetLogger()

	err := os.MkdirAll(c.ModelsPath, os.ModePerm)
	if err != nil {
		logger.Fatal("erro ao criar a pasta %s: %v", c.ModelsPath, err)
	}

	err = os.MkdirAll(c.TablesPath, os.ModePerm)
	if err != nil {
		logger.Fatal("erro ao criar a pasta %s: %v", c.TablesPath, err)

	}
	logger.Info("Project Created xxx")
	return nil
}

func (c *Config) CreateExampleProject() {
	logger := logger.GetLogger()

	modelsExample := make(map[string]string)
	tablesExample := make(map[string]string)

	tablesExample["todo.sql"] = `
-- sql2api.description: Criação da tabela todos

CREATE TABLE IF NOT EXISTS todos (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task TEXT NOT NULL,
	completed INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
	`

	modelsExample["create_one.sql"] = `
-- sql2api.name: create_task
-- sql2api.description: usecase for task creation in context of TODO app
-- sql2api.http: POST
-- sql2api.async: yes
-- sql2api.parameters: magic

INSERT INTO todos (task, completed) VALUES ('Comprar leite', 1);
	`

	modelsExample["delete_one.sql"] = `
-- sql2api.name: del_task_by_id
-- sql2api.description: usecase for task deleting in context of TODO app
-- sql2api.http: DELETE
-- sql2api.async: yes
-- sql2api.parameters: id

DELETE FROM todos WHERE id = 1;
		`

	modelsExample["read_all.sql"] = `
-- sql2api.name: get_all_tasks
-- sql2api.description: usecase for task reading in context of TODO app
-- sql2api.http: GET
-- sql2api.async: yes
-- sql2api.parameters: max=10

SELECT * FROM todos;
	`

	modelsExample["read_one.sql"] = `
-- sql2api.name: get_task_by_id
-- sql2api.description: usecase for task reading in context of TODO app
-- sql2api.http: GET
-- sql2api.async: yes
-- sql2api.parameters: id

SELECT * FROM todos WHERE id = 1;
	`

	modelsExample["update_one.sql"] = `
-- sql2api.name: edit_task_by_id
-- sql2api.description: usecase for task updating in context of TODO app
-- sql2api.http: PUT
-- sql2api.async: yes
-- sql2api.parameters: auto

UPDATE todos
SET task = 'Comprar leite e pão', completed = TRUE
WHERE id = 1;
	`

	for filename, sqlContent := range modelsExample {
		filePath := filepath.Join(c.ModelsPath, filename)
		file, err := os.Create(filePath)
		if err != nil {
			logger.Fatalf("Error creating SQL file %s: %v", filePath, err)
		}
		defer file.Close()

		_, err = file.WriteString(sqlContent)
		if err != nil {
			logger.Fatalf("Error writing to SQL file %s: %v", filePath, err)
		}
		logger.Printf("Successfully created file: %s", filePath)
	}

	for filename, sqlContent := range tablesExample {
		filePath := filepath.Join(c.TablesPath, filename)
		file, err := os.Create(filePath)
		if err != nil {
			logger.Fatalf("Error creating SQL file %s: %v", filePath, err)
		}
		defer file.Close()

		_, err = file.WriteString(sqlContent)
		if err != nil {
			logger.Fatalf("Error writing to SQL file %s: %v", filePath, err)
		}
		logger.Printf("Successfully created file: %s", filePath)
	}
}

func (c *Config) SetupModels(models map[string]string) (map[string]models.ModelConfig, error) {
	// Initialize the modelConfig map
	modelConfig := make(map[string]models.ModelConfig)
	// Populate the modelConfig map
	for fileName, content := range models {
		mConfig, err := c.parseModelSetup(content)
		if err != nil {
			return nil, err
		}
		modelConfig[fileName] = mConfig
	}
	return modelConfig, nil
}

func (c *Config) parseModelSetup(sql string) (models.ModelConfig, error) {
	var modelConfig models.ModelConfig
	scanner := bufio.NewScanner(strings.NewReader(sql))
	var query string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// ignore sql
		if len(line) == 0 || !strings.Contains(line, "sql2api") {
			query += line
			continue
		}

		// Remover o prefixo "-- sql2api" e fazer um split
		trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "-- sql2api."))

		switch {
		case strings.HasPrefix(trimmedLine, "description:"):
			modelConfig.Description = strings.TrimPrefix(trimmedLine, "description:")
		case strings.HasPrefix(trimmedLine, "http:"):
			modelConfig.HTTPMethod = strings.TrimPrefix(trimmedLine, "http:")
		case strings.HasPrefix(trimmedLine, "async:"):
			asyncValue := strings.TrimPrefix(trimmedLine, "async:")
			if asyncValue == "yes" {
				modelConfig.Async = true
			} else {
				modelConfig.Async = false
			}
		case strings.HasPrefix(trimmedLine, "parameters:"):
			modelConfig.Parameters = strings.TrimPrefix(trimmedLine, "parameters:")
		}
	}

	modelConfig.Query = query
	if err := scanner.Err(); err != nil {
		return models.ModelConfig{}, err
	}
	return modelConfig, nil
}
