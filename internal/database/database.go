package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Database interface {
	Connect() error
	Execute(sqlQuery string) error
}

type SQLiteDatabase struct {
	DatabaseFile string
	Conn         *sql.DB
}

func (db *SQLiteDatabase) Connect() error {
	conn, err := sql.Open("sqlite3", db.DatabaseFile)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco de dados: %w", err)
	}
	log.Println("Database", db.DatabaseFile, "criado com sucesso!")
	db.Conn = conn
	return nil
}

func (db *SQLiteDatabase) Execute(sqlQuery string) error {
	_, err := db.Conn.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("erro ao executar a consulta: %w", err)
	}
	return nil
}
