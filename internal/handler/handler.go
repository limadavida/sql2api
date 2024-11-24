package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/limadavida/sql2api/internal/config"
	"github.com/limadavida/sql2api/internal/database"
	"github.com/limadavida/sql2api/internal/utils"
)

type Handler interface {
	Router()
}

type Routes struct {
	cfg         config.Config
	models      utils.SqlNamed
	dbConnector database.Database
	ginProvider *gin.Engine
}

func NewHandler(cfg config.Config) Handler {
	models := config.ValidateSqlModels(cfg.RootDir + "/models")

	db, err := database.NewDatabase(cfg.Databases.Type, cfg.Databases.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &Routes{
		cfg:         cfg,
		models:      models,
		dbConnector: db,
		ginProvider: gin.Default(),
	}
}

func (r *Routes) Router() {

	basicRoute := "/" + r.cfg.Project

	r.ginProvider.POST(basicRoute, r.executeQuery(r.cfg.RootDir+"/models/"+r.cfg.Models.POST.File[0]))
	r.ginProvider.GET(basicRoute, r.executeQuery(r.cfg.RootDir+"/models/"+r.cfg.Models.POST.File[0]))
	r.ginProvider.PUT(basicRoute, r.executeQuery(r.cfg.RootDir+"/models/"+r.cfg.Models.PUT.File[0]))
	r.ginProvider.DELETE(basicRoute, r.executeQuery(r.cfg.RootDir+"/models/"+r.cfg.Models.DELETE.File[0]))

	port := r.cfg.Servers[0]
	r.ginProvider.Run(fmt.Sprintf(":%d", port))
}

func (r *Routes) executeQuery(modelFile string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		sqlQuery := utils.ReadFile(modelFile)

		log.Println("sqlQuery\n", sqlQuery)

		err := r.dbConnector.Execute(sqlQuery)
		if err != nil {
			c.String(500, modelFile)
		}
		c.String(200, modelFile)
	}
	return gin.HandlerFunc(fn)

}
