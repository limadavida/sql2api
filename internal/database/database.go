package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Database interface {
	Execute(sqlQuery string) error
}

type GormDatabase struct {
	DB *gorm.DB
}

func NewDatabase(dbType, dsn string) (Database, error) {

	var db *gorm.DB
	var err error

	switch dbType {
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgresql": // db, err = database.NewDatabase("postgres", "user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("tipo de banco de dados desconhecido: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conectado ao banco de dados", dbType, "com sucesso!")

	return &GormDatabase{DB: db}, nil

}

func (gdb *GormDatabase) Execute(sqlQuery string) error {
	result := gdb.DB.Exec(sqlQuery)
	if result.Error != nil {
		return fmt.Errorf("erro ao executar consulta: %w", result.Error)
	}
	return nil
}
