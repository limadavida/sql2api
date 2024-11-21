package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/limadavida/sql2api/internal/config"
	"github.com/limadavida/sql2api/internal/database"
	"github.com/limadavida/sql2api/internal/utils"
)

type Handler interface {
	Post() gin.HandlerFunc
	Get() gin.HandlerFunc
	Put() gin.HandlerFunc
	Del() gin.HandlerFunc
}

type Routes struct {
	cfg       config.Config
	models    utils.SqlNamed
	dbService *database.SQLiteDatabase
}

func NewHandler(cfg config.Config) Handler {
	models := config.ValidateSqlModels(cfg.RootDir + "/models")
	sqliteDB := &database.SQLiteDatabase{DatabaseFile: config.ConfigData.Databases.Name}

	return &Routes{cfg: cfg, models: models, dbService: sqliteDB}
}

func (r *Routes) Post() gin.HandlerFunc {
	err := r.dbService.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer r.dbService.Conn.Close()

	sqlNames := utils.RemoveExtFromList(r.cfg.Models.POST.File, ".sql")
	for _, sqlName := range sqlNames {
		sqlQuery := r.models[sqlName]
		log.Println(sqlQuery)
		err := r.dbService.Execute(sqlQuery)
		if err != nil {
			return func(c *gin.Context) {
				c.JSON(500, gin.H{"error": "Internal Server Error"})
			}
		}
	}

	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST request for createTask", "models": r.cfg.Models.POST.File})
	}
}

func (r *Routes) Get() gin.HandlerFunc {
	err := r.dbService.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer r.dbService.Conn.Close()

	sqlNames := utils.RemoveExtFromList(r.cfg.Models.POST.File, ".sql")
	for _, sqlName := range sqlNames {
		sqlQuery := r.models[sqlName]
		log.Println(sqlQuery)
		err := r.dbService.Execute(sqlQuery)
		if err != nil {
			return func(c *gin.Context) {
				c.JSON(500, gin.H{"error": "Internal Server Error"})
			}
		}
	}
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GET request for readTask", "files": r.cfg.Models.GET.File})
	}
}

func (r *Routes) Put() gin.HandlerFunc {
	err := r.dbService.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer r.dbService.Conn.Close()

	sqlNames := utils.RemoveExtFromList(r.cfg.Models.POST.File, ".sql")
	for _, sqlName := range sqlNames {
		sqlQuery := r.models[sqlName]
		log.Println(sqlQuery)
		err := r.dbService.Execute(sqlQuery)
		if err != nil {
			return func(c *gin.Context) {
				c.JSON(500, gin.H{"error": "Internal Server Error"})
			}
		}
	}
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "PUT request for updateTask", "files": r.cfg.Models.PUT})
	}
}

func (r *Routes) Del() gin.HandlerFunc {
	err := r.dbService.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer r.dbService.Conn.Close()

	sqlNames := utils.RemoveExtFromList(r.cfg.Models.POST.File, ".sql")
	for _, sqlName := range sqlNames {
		sqlQuery := r.models[sqlName]
		log.Println(sqlQuery)
		err := r.dbService.Execute(sqlQuery)
		if err != nil {
			return func(c *gin.Context) {
				c.JSON(500, gin.H{"error": "Internal Server Error"})
			}
		}
	}
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "DELETE request for deleteTask", "files": r.cfg.Models.DELETE})
	}
}
