package main

import (
	"fmt"

	"github.com/limadavida/sql2api/internal/config"
	_ "github.com/limadavida/sql2api/internal/config"
	"github.com/limadavida/sql2api/internal/handler"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("[SQL2API] Starting...")

	config.CreateTables()
	handler.Router()
}
