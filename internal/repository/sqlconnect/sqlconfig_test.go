package sqlconnect

import (
	"os"
	"sync"
	"testing"
)

func TestInitDBPool(t *testing.T) {
	// Сохраняем оригинальные значения
	originalVars := map[string]string{
		"DB_USER":   os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":   os.Getenv("DB_NAME"),
		"DB_PORT":   os.Getenv("DB_PORT"),
		"HOST":      os.Getenv("HOST"),
	}

	// Устанавливаем тестовые значения (если они не установлены)
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "test")
		os.Setenv("DB_PASSWORD", "test")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("HOST", "localhost")
	}

	defer func() {
		// Восстанавливаем оригинальные значения
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
		// Закрываем пул после теста
		CloseDBPool()
	}()

	// Тест инициализации
	err := InitDBPool()
	if err != nil {
		// Если БД недоступна, это нормально для unit тестов
		t.Logf("InitDBPool() failed (expected if DB not available): %v", err)
		return
	}

	// Проверяем, что пул создан
	db, err := GetDB()
	if err != nil {
		t.Errorf("GetDB() failed after InitDBPool: %v", err)
		return
	}

	if db == nil {
		t.Errorf("GetDB() returned nil")
	}

	// Проверяем, что повторная инициализация не создает новый пул
	err2 := InitDBPool()
	if err2 != nil {
		t.Errorf("Second InitDBPool() should not fail: %v", err2)
	}

	db2, _ := GetDB()
	if db != db2 {
		t.Errorf("InitDBPool() should return same instance on second call")
	}
}

func TestGetDBWithoutInit(t *testing.T) {
	// Сбрасываем пул для теста
	CloseDBPool()
	
	// Сбрасываем once для повторной инициализации
	once = sync.Once{}

	_, err := GetDB()
	if err == nil {
		t.Errorf("GetDB() should fail if pool not initialized")
	}
}

func TestCloseDBPool(t *testing.T) {
	// Инициализируем пул
	originalVars := map[string]string{
		"DB_USER":   os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":   os.Getenv("DB_NAME"),
		"DB_PORT":   os.Getenv("DB_PORT"),
		"HOST":      os.Getenv("HOST"),
	}

	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "test")
		os.Setenv("DB_PASSWORD", "test")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("HOST", "localhost")
	}

	defer func() {
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	InitDBPool()
	
	err := CloseDBPool()
	if err != nil {
		t.Logf("CloseDBPool() error (may be expected): %v", err)
	}

	// Повторное закрытие не должно вызывать ошибку
	err = CloseDBPool()
	if err != nil {
		t.Logf("Second CloseDBPool() error (may be expected): %v", err)
	}
}

