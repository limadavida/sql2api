package main

import (
	"encoding/json"

	"github.com/limadavida/sql2api/internal/cli"
	"github.com/limadavida/sql2api/internal/config"
	_ "github.com/limadavida/sql2api/internal/config"
	"github.com/limadavida/sql2api/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := logger.GetLogger()
	logger.Info("[SQL2API] Starting...")

	appConfig := cli.RunSetup()

	logger.Info("Application Settings:")
	data, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		logger.Error("Erro ao gerar JSON:", err)
		return
	}
	logger.Info(string(data))

	config.CreateProjectPath(appConfig.Project)

	logger.Info("Add your sql tables at /tables. cli com botao de ok!")
	config.CreateTables(appConfig)

	logger.Info("Add your sql  at /models. cli com botao de ok!")

	//handler := handler.NewHandler(*config.ConfigData)
	//handler.Router()

}
