package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/limadavida/sql2api/internal/database"
	"github.com/limadavida/sql2api/internal/logger"
	"github.com/limadavida/sql2api/internal/models"
	"github.com/limadavida/sql2api/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Project   string `mapstructure:"project"`
	Version   string `mapstructure:"version"`
	Author    string `mapstructure:"author"`
	Servers   []int  `mapstructure:"servers"`
	Type      string `mapstructure:"type"`
	RootDir   string `mapstructure:"rootDir"`
	Databases struct {
		Type string `mapstructure:"type"`
		Name string `mapstructure:"name"`
	} `mapstructure:"databases"`

	// Models agora usa as chaves HTTP como POST, GET, PUT, DELETE
	Models struct {
		POST   ModelConfig `mapstructure:"POST"`
		GET    ModelConfig `mapstructure:"GET"`
		PUT    ModelConfig `mapstructure:"PUT"`
		DELETE ModelConfig `mapstructure:"DELETE"`
	} `mapstructure:"models"`
}

// ModelConfig agora contém os arquivos e a ação HTTP
type ModelConfig struct {
	File []string `mapstructure:"file"` // Lista de arquivos SQL
}

func loadConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("erro ao unmarshaling o arquivo de configuração: %w", err)
	}

	return &config, nil
}

func validateSqlFiles(directory string) map[string]string {
	sqlMap := make(map[string]string)
	sqlFiles, _ := utils.ListFiles(directory, ".sql")

	if len(sqlFiles) == 0 {
		log.Println("Nao foi possível encontrar modelos. Verifique se há arquivos SQL válidos em", directory)
	}

	var fails []string
	for _, sqlFile := range sqlFiles {
		sqlBytes, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			log.Fatalf("Erro ao ler o arquivo SQL: %v\n", err)
		}
		sqlFileName := strings.Replace(filepath.Base(sqlFile), ".sql", "", 1)
		sqlQuery := string(sqlBytes)

		if utils.IsEmpty(sqlQuery) {
			fails = append(fails, sqlFileName)
		} else {
			sqlMap[sqlFileName] = sqlQuery
			log.Println("Model", sqlFileName, "encontrado em", sqlFile)
		}
	}

	if len(fails) > 0 {
		log.Printf("The following SQL files are empty: %v\n", fails)
	}

	return sqlMap
}

func ValidateSqlModels(directory string) utils.SqlNamed {
	return validateSqlFiles(directory)
}

func validateSqlTables(directory string) utils.SqlNamed {
	return validateSqlFiles(directory)
}

func CreateTables(cfg models.ProjectConfig) {
	logger := logger.GetLogger()
	tables := validateSqlTables(cfg.RootDir + "/tables")

	db, err := database.NewDatabase(cfg.Database.Type, cfg.Database.Name)
	if err != nil {
		logger.Fatal(err)
	}

	for tableName, tableSql := range tables {
		err = db.Execute(tableSql)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Tabela", tableName, "criada com sucesso!")
	}
}

func CreateProjectPath(projectName string) error {
	logger := logger.GetLogger()

	projectBasePath := filepath.Join(".", projectName)
	modelsPath := filepath.Join(projectBasePath, "models")
	tablesPath := filepath.Join(projectBasePath, "tables")

	err := os.MkdirAll(modelsPath, os.ModePerm)
	if err != nil {
		logger.Fatal("erro ao criar a pasta %s: %v", projectName, err)
	}

	err = os.MkdirAll(modelsPath, os.ModePerm)
	if err != nil {
		logger.Fatal("erro ao criar a pasta 'models': %v", err)
	}

	err = os.MkdirAll(tablesPath, os.ModePerm)
	if err != nil {
		logger.Fatal("erro ao criar a pasta 'tables': %v", err)
	}

	logger.Info("Project Created at", projectBasePath)
	return nil
}


