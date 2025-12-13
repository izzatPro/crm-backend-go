package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbPool *sql.DB
	once   sync.Once
	mu     sync.RWMutex
)

// InitDBPool инициализирует глобальный пул соединений с БД
func InitDBPool() error {
	var initErr error
	once.Do(func() {
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		dbport := os.Getenv("DB_PORT")
		host := os.Getenv("HOST")

		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, dbport, dbname)
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			initErr = err
			return
		}

		// Настройка пула соединений
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(0) // Соединения не закрываются по таймауту

		// Проверка соединения
		if err := db.Ping(); err != nil {
			initErr = err
			return
		}

		dbPool = db
	})
	return initErr
}

// GetDB возвращает глобальный пул соединений
func GetDB() (*sql.DB, error) {
	mu.RLock()
	defer mu.RUnlock()
	if dbPool == nil {
		return nil, fmt.Errorf("database pool not initialized, call InitDBPool first")
	}
	return dbPool, nil
}

// CloseDBPool закрывает пул соединений
func CloseDBPool() error {
	mu.Lock()
	defer mu.Unlock()
	if dbPool != nil {
		return dbPool.Close()
	}
	return nil
}

// ConnectDb оставлен для обратной совместимости, но теперь использует пул
func ConnectDb() (*sql.DB, error) {
	return GetDB()
}
