// Package database отвечает за создание БД и CRUD операции
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"wb_project_0/config"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

type Database struct {
	conn   *sql.DB
	logger *log.Logger
}

func NewDB(db *sql.DB) *Database {
	return &Database{
		conn:   db,
		logger: log.New(os.Stdout, "[DB] ", log.LstdFlags|log.Lshortfile),
	}
}

// InitDB инициализирует подключение к базе данных
func InitDB() (*Database, error) {
	connStr, err := config.GetDBConf()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB config: %w", err)
	}

	db, err := sql.Open("postgres", connStr) // Явно указать драйвер
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	database := NewDB(db)
	database.logger.Println("Initializing database connection")
	// Проверка подключения
	if err := db.Ping(); err != nil {
		database.logger.Printf("Connection ping failed: %v", err)
		database.conn.Close() // Закрыть соединение при ошибке
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	if err := database.CreateIfNotExists(); err != nil {
		database.logger.Printf("Failed to initialize tables: %v", err)
		database.conn.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	database.logger.Println("Database initialized successfully")
	return database, nil
}

// CheckAllTablesExist проверяет существование всех таблиц
func (db *Database) CheckAllTablesExist() (bool, error) {
	query := `
	SELECT COUNT(*) = 4 AS all_tables_exist
	FROM information_schema.tables 
	WHERE table_schema = 'public' 
	AND table_name IN ('orders', 'deliveries', 'payments', 'items')`

	var allExist bool
	err := db.conn.QueryRow(query).Scan(&allExist)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}
	return allExist, nil
}

// CreateIfNotExists создает таблицы, если они не существуют
func (db *Database) CreateIfNotExists() error {
	exists, err := db.CheckAllTablesExist()
	if err != nil {
		return fmt.Errorf("failed to check tables existence: %w", err)
	}

	if !exists {
		if err := db.CreateTables(); err != nil {
			return fmt.Errorf("failed to create tables: %w", err)
		}
	}
	return nil
}

// CreateTables создает все таблицы
func (db *Database) CreateTables() error {
	migration, err := os.ReadFile("/app/migrations/002_init_db.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := db.conn.Exec(string(migration)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	return nil
}

// DeleteTables удаляет все таблицы
func (db *Database) DeleteTables() error {
	migration, err := os.ReadFile("/app/migrations/001_delete_db.sql")
	if err != nil {
		return fmt.Errorf("failed to read delete script: %w", err)
	}

	if _, err := db.conn.Exec(string(migration)); err != nil {
		return fmt.Errorf("failed to execute delete script: %w", err)
	}
	return nil
}

// Close закрывает соединение с базой данных
func (db *Database) Close() error {
	if db.conn == nil {
		return nil
	}
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}
	return nil
}
