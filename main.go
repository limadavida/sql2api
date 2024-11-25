package main

import (
	"encoding/json"

	"github.com/limadavida/sql2api/internal/cli"
	"github.com/limadavida/sql2api/internal/config"
	_ "github.com/limadavida/sql2api/internal/config"
	"github.com/limadavida/sql2api/internal/database"
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

	db, err := database.NewDatabase(appConfig.Database.Type, appConfig.Database.Name)
	if err != nil {
		logger.Fatal(err)
	}

	projectSetup := config.NewConfig(appConfig.RootDir, appConfig.Project)

	projectSetup.CreateProjectPath()

	logger.Info("Add your sql tables at /tables. cli com botao de ok!")
	logger.Info("Add your sql  at /models. cli com botao de ok!")

	projectSetup.CreateExampleProject() //TODO: if setup is empty

	models, tables, err := projectSetup.GetProject()
	if err != nil {
		logger.Fatal(err)
	}

	err = projectSetup.CreateTables(db, tables)
	if err != nil {
		logger.Fatal(err)
	}

	projectSetup.SetupModels(models)

	//handler := handler.NewHandler(*config.ConfigData)
	//handler.Router()

}
