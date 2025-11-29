package test

import (
	"fmt"
	"go-graphql/internal/config"
	"strings"
	"testing"
)

// TestValidateConfigValidComplete tests valid complete configuration
func TestValidateConfigValidComplete(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://user:pass@localhost:5432/database?sslmode=disable",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err != nil {
		t.Fatalf("❌ Expected no error, got: %v", err)
	}

	t.Log("✅ TestValidateConfigValidComplete passed")
}

// TestValidateConfigInvalidHTTPPort tests invalid HTTP port
func TestValidateConfigInvalidHTTPPort(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    99999,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid HTTP_PORT") {
		t.Fatalf("❌ Expected error containing 'invalid HTTP_PORT', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidHTTPPort passed")
}

// TestValidateConfigEmptyAddress tests empty HTTP address
func TestValidateConfigEmptyAddress(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "HTTP_ADDRESS is empty") {
		t.Fatalf("❌ Expected error containing 'HTTP_ADDRESS is empty', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigEmptyAddress passed")
}

// TestValidateConfigInvalidEnvironment tests invalid environment
func TestValidateConfigInvalidEnvironment(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "staging",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid ENV") {
		t.Fatalf("❌ Expected error containing 'invalid ENV', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidEnvironment passed")
}

// TestValidateConfigEmptyEnvironment tests empty environment
func TestValidateConfigEmptyEnvironment(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "ENV is empty") {
		t.Fatalf("❌ Expected error containing 'ENV is empty', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigEmptyEnvironment passed")
}

// TestValidateConfigEmptyDatabaseDSN tests empty database DSN
func TestValidateConfigEmptyDatabaseDSN(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DATABASE_DSN is empty") {
		t.Fatalf("❌ Expected error containing 'DATABASE_DSN is empty', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigEmptyDatabaseDSN passed")
}

// TestValidateConfigInvalidDatabaseDSN tests invalid database DSN format
func TestValidateConfigInvalidDatabaseDSN(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "localhost:5432/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid DATABASE_DSN format") {
		t.Fatalf("❌ Expected error containing 'invalid DATABASE_DSN format', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidDatabaseDSN passed")
}

// TestValidateConfigEmptyRedisDSN tests empty Redis DSN
func TestValidateConfigEmptyRedisDSN(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "REDIS_DSN is empty") {
		t.Fatalf("❌ Expected error containing 'REDIS_DSN is empty', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigEmptyRedisDSN passed")
}

// TestValidateConfigInvalidRedisDSN tests invalid Redis DSN format
func TestValidateConfigInvalidRedisDSN(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid REDIS_DSN format") {
		t.Fatalf("❌ Expected error containing 'invalid REDIS_DSN format', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidRedisDSN passed")
}

// TestValidateConfigInvalidRedisDB tests invalid Redis DB
func TestValidateConfigInvalidRedisDB(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         20,
			Prefix:     "go-graphql",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid REDIS_DB") {
		t.Fatalf("❌ Expected error containing 'invalid REDIS_DB', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidRedisDB passed")
}

// TestValidateConfigEmptyRedisPrefix tests empty Redis prefix
func TestValidateConfigEmptyRedisPrefix(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "REDIS_PREFIX is empty") {
		t.Fatalf("❌ Expected error containing 'REDIS_PREFIX is empty', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigEmptyRedisPrefix passed")
}

// TestValidateConfigInvalidRedisPrefix tests Redis prefix with spaces
func TestValidateConfigInvalidRedisPrefix(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go simple",
			DefaultTTL: 5,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid REDIS_PREFIX") {
		t.Fatalf("❌ Expected error containing 'invalid REDIS_PREFIX', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidRedisPrefix passed")
}

// TestValidateConfigInvalidRedisTTL tests invalid Redis TTL
func TestValidateConfigInvalidRedisTTL(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    4000,
		HTTPAddress: "127.0.0.1",
		ENV:         "development",
		Database: config.DatabaseCfg{
			DSN: "postgresql://localhost/db",
		},
		Redis: config.RedisCfg{
			DSN:        "localhost:6379",
			DB:         0,
			Prefix:     "go-graphql",
			DefaultTTL: 0,
		},
	}

	// Test
	err := config.ValidateConfig(cfg)

	// Assert
	if err == nil {
		t.Fatalf("❌ Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid REDIS_DEFAULT_TTL") {
		t.Fatalf("❌ Expected error containing 'invalid REDIS_DEFAULT_TTL', got: %v", err)
	}

	fmt.Println(err.Error())
	t.Log("✅ TestValidateConfigInvalidRedisTTL passed")
}
