package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/limadavida/sql2api/internal/database"
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

func CreateTables() {
	tables := validateSqlTables("examples/TodoExample/tables")
	sqliteDB := &database.SQLiteDatabase{DatabaseFile: ConfigData.Databases.Name}

	err := sqliteDB.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDB.Conn.Close()

	for tableName, tableSql := range tables {
		err = sqliteDB.Execute(tableSql)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Tabela", tableName, "criada com sucesso!")
	}
}
