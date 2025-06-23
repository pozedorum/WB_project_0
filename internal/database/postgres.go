// Package database is responcible for creating database and CRUD operations
package database

import (
	"database/sql"
	"fmt"
	"os"
	"wb_project_0/config"
	"wb_project_0/internal/models"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

func InitDB() (*Database, error) {
	url, err := config.GetDBConf()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{conn: db}, nil
}

func (db *Database) CreateIfNotExisists() error {
	exists, err := db.IsTableExists()
	if err != nil {
		return err
	}
	if !exists {
		err = db.CreateTable()
		return err
	}
	return nil
}

func (db *Database) CheckAllTablesExist() (bool, error) {
	query := `
    SELECT COUNT(*) = 4 AS all_tables_exist
    FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name IN ('orders', 'deliveries', 'payments', 'items')
    `

	var allExist bool
	err := db.conn.QueryRow(query).Scan(&allExist)
	return allExist, err
}

func (db *Database) CreateTables() error {
	migration, err := os.ReadFile("migrations/init_db.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migration))
	return err
}

func (db *Database) DeleteTables() error {
	migration, err := os.ReadFile("migrations/delete_db.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migration))
	return err
}

func (db *Database) Close() error {
	if db.conn == nil {
		return nil // или возвращаем ошибку, если это не ожидаемое состояние
	}
	return db.conn.Close()
}
