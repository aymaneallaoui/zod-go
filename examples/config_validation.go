package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// Configuration schemas
var (
	// Database configuration schema
	databaseConfigSchema = validators.Object(map[string]zod.Schema{
		"host": validators.String().
			Min(1).
			Required().
			WithMessage("required", "Database host is required"),
		"port": validators.Number().
			Integer().
			Min(1).
			Max(65535).
			Default(5432),
		"name": validators.String().
			Min(1).
			Required().
			WithMessage("required", "Database name is required"),
		"username": validators.String().
			Min(1).
			Required(),
		"password": validators.String().
			Min(1).
			Required(),
		"sslMode": validators.String().
			Pattern(`^(disable|require|verify-ca|verify-full)$`).
			Default("require"),
		"maxConnections": validators.Number().
			Integer().
			Min(1).
			Max(1000).
			Default(25),
		"connectionTimeout": validators.Number().
			Integer().
			Min(1).
			Max(300).
			Default(30),
	})

	// Redis configuration schema
	redisConfigSchema = validators.Object(map[string]zod.Schema{
		"host":        validators.String().Min(1).Default("localhost"),
		"port":        validators.Number().Integer().Min(1).Max(65535).Default(6379),
		"password":    validators.String().Optional(),
		"database":    validators.Number().Integer().Min(0).Max(15).Default(0),
		"maxRetries":  validators.Number().Integer().Min(0).Default(3),
		"dialTimeout": validators.Number().Integer().Min(1).Default(5),
	})

	// Server configuration schema
	serverConfigSchema = validators.Object(map[string]zod.Schema{
		"host": validators.String().Default("0.0.0.0"),
		"port": validators.Number().
			Integer().
			Min(1).
			Max(65535).
			Default(8080),
		"readTimeout": validators.Number().
			Integer().
			Min(1).
			Max(300).
			Default(30),
		"writeTimeout": validators.Number().
			Integer().
			Min(1).
			Max(300).
			Default(30),
		"maxHeaderBytes": validators.Number().
			Integer().
			Min(1024).
			Default(1048576), // 1MB
		"tls": validators.Object(map[string]zod.Schema{
			"enabled":  validators.Bool().Default(false),
			"certFile": validators.String().Optional(),
			"keyFile":  validators.String().Optional(),
		}).Optional(),
	})

	// Logging configuration schema
	loggingConfigSchema = validators.Object(map[string]zod.Schema{
		"level": validators.String().
			Pattern(`^(debug|info|warn|error|fatal)$`).
			Default("info"),
		"format": validators.String().
			Pattern(`^(json|text)$`).
			Default("json"),
		"output": validators.String().
			Pattern(`^(stdout|stderr|file)$`).
			Default("stdout"),
		"file": validators.Object(map[string]zod.Schema{
			"path":       validators.String().Optional(),
			"maxSize":    validators.Number().Integer().Min(1).Default(100), // MB
			"maxBackups": validators.Number().Integer().Min(0).Default(3),
			"maxAge":     validators.Number().Integer().Min(0).Default(28), // days
			"compress":   validators.Bool().Default(true),
		}).Optional(),
	})

	// Application configuration schema
	appConfigSchema = validators.Object(map[string]zod.Schema{
		"environment": validators.String().
			Pattern(`^(development|staging|production)$`).
			Default("development"),
		"debug": validators.Bool().Default(false),
		"appName": validators.String().
			Min(1).
			Required().
			WithMessage("required", "Application name is required"),
		"version": validators.String().
			Pattern(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$`).
			Required(),
		"secretKey": validators.String().
			Min(32).
			Required().
			WithMessage("minLength", "Secret key must be at least 32 characters"),
		"cors": validators.Object(map[string]zod.Schema{
			"allowedOrigins": validators.Array(validators.String().URL().Required()).
				Min(1).
				Default([]interface{}{"http://localhost:3000"}),
			"allowedMethods": validators.Array(validators.String().Required()).
				Default([]interface{}{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			"allowedHeaders": validators.Array(validators.String().Required()).
				Default([]interface{}{"Content-Type", "Authorization"}),
			"allowCredentials": validators.Bool().Default(false),
		}).Optional(),
		"rateLimit": validators.Object(map[string]zod.Schema{
			"enabled":        validators.Bool().Default(true),
			"requestsPerMin": validators.Number().Integer().Min(1).Default(100),
			"burst":          validators.Number().Integer().Min(1).Default(20),
		}).Optional(),
	})

	// Complete application configuration schema
	configSchema = validators.Object(map[string]zod.Schema{
		"app":      appConfigSchema.Required(),
		"server":   serverConfigSchema.Required(),
		"database": databaseConfigSchema.Required(),
		"redis":    redisConfigSchema.Optional(),
		"logging":  loggingConfigSchema.Required(),
	}).Strict()
)

// Configuration struct
type Config struct {
	App      AppConfig      `json:"app"`
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    *RedisConfig   `json:"redis,omitempty"`
	Logging  LoggingConfig  `json:"logging"`
}

type AppConfig struct {
	Environment string      `json:"environment"`
	Debug       bool        `json:"debug"`
	AppName     string      `json:"appName"`
	Version     string      `json:"version"`
	SecretKey   string      `json:"secretKey"`
	CORS        *CORSConfig `json:"cors,omitempty"`
	RateLimit   *RateLimit  `json:"rateLimit,omitempty"`
}

type CORSConfig struct {
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedHeaders   []string `json:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
}

type RateLimit struct {
	Enabled        bool `json:"enabled"`
	RequestsPerMin int  `json:"requestsPerMin"`
	Burst          int  `json:"burst"`
}

type ServerConfig struct {
	Host           string     `json:"host"`
	Port           int        `json:"port"`
	ReadTimeout    int        `json:"readTimeout"`
	WriteTimeout   int        `json:"writeTimeout"`
	MaxHeaderBytes int        `json:"maxHeaderBytes"`
	TLS            *TLSConfig `json:"tls,omitempty"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
}

type DatabaseConfig struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Name              string `json:"name"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	SSLMode           string `json:"sslMode"`
	MaxConnections    int    `json:"maxConnections"`
	ConnectionTimeout int    `json:"connectionTimeout"`
}

type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password,omitempty"`
	Database    int    `json:"database"`
	MaxRetries  int    `json:"maxRetries"`
	DialTimeout int    `json:"dialTimeout"`
}

type LoggingConfig struct {
	Level  string   `json:"level"`
	Format string   `json:"format"`
	Output string   `json:"output"`
	File   *FileLog `json:"file,omitempty"`
}

type FileLog struct {
	Path       string `json:"path,omitempty"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

// LoadConfigFromJSON loads and validates configuration from JSON file
func LoadConfigFromJSON(filename string) (*Config, error) {
	fmt.Printf("Loading configuration from %s...\n", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configData map[string]interface{}
	if err := json.Unmarshal(data, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate configuration
	if err := configSchema.Validate(configData); err != nil {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			return nil, fmt.Errorf("configuration validation failed:\n%s", validationErr.ErrorJSON())
		}
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Convert to struct
	var config Config
	configJSON, _ := json.Marshal(configData)
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, fmt.Errorf("failed to convert to config struct: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads and validates configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	fmt.Println("Loading configuration from environment variables...")

	configData := map[string]interface{}{
		"app": map[string]interface{}{
			"environment": getEnvOrDefault("APP_ENV", "development"),
			"debug":       getEnvBoolOrDefault("APP_DEBUG", false),
			"appName":     getEnvOrDefault("APP_NAME", ""),
			"version":     getEnvOrDefault("APP_VERSION", ""),
			"secretKey":   getEnvOrDefault("APP_SECRET_KEY", ""),
		},
		"server": map[string]interface{}{
			"host":           getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
			"port":           getEnvIntOrDefault("SERVER_PORT", 8080),
			"readTimeout":    getEnvIntOrDefault("SERVER_READ_TIMEOUT", 30),
			"writeTimeout":   getEnvIntOrDefault("SERVER_WRITE_TIMEOUT", 30),
			"maxHeaderBytes": getEnvIntOrDefault("SERVER_MAX_HEADER_BYTES", 1048576),
		},
		"database": map[string]interface{}{
			"host":              getEnvOrDefault("DB_HOST", ""),
			"port":              getEnvIntOrDefault("DB_PORT", 5432),
			"name":              getEnvOrDefault("DB_NAME", ""),
			"username":          getEnvOrDefault("DB_USERNAME", ""),
			"password":          getEnvOrDefault("DB_PASSWORD", ""),
			"sslMode":           getEnvOrDefault("DB_SSL_MODE", "require"),
			"maxConnections":    getEnvIntOrDefault("DB_MAX_CONNECTIONS", 25),
			"connectionTimeout": getEnvIntOrDefault("DB_CONNECTION_TIMEOUT", 30),
		},
		"logging": map[string]interface{}{
			"level":  getEnvOrDefault("LOG_LEVEL", "info"),
			"format": getEnvOrDefault("LOG_FORMAT", "json"),
			"output": getEnvOrDefault("LOG_OUTPUT", "stdout"),
		},
	}

	// Add Redis configuration if provided
	if redisHost := getEnvOrDefault("REDIS_HOST", ""); redisHost != "" {
		configData["redis"] = map[string]interface{}{
			"host":        redisHost,
			"port":        getEnvIntOrDefault("REDIS_PORT", 6379),
			"password":    getEnvOrDefault("REDIS_PASSWORD", ""),
			"database":    getEnvIntOrDefault("REDIS_DATABASE", 0),
			"maxRetries":  getEnvIntOrDefault("REDIS_MAX_RETRIES", 3),
			"dialTimeout": getEnvIntOrDefault("REDIS_DIAL_TIMEOUT", 5),
		}
	}

	// Validate configuration
	if err := configSchema.Validate(configData); err != nil {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			return nil, fmt.Errorf("environment configuration validation failed:\n%s", validationErr.ErrorJSON())
		}
		return nil, fmt.Errorf("environment configuration validation failed: %w", err)
	}

	// Convert to struct
	var config Config
	configJSON, _ := json.Marshal(configData)
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, fmt.Errorf("failed to convert to config struct: %w", err)
	}

	return &config, nil
}

// Environment variable helpers
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func configMain() {
	fmt.Println("=== Zod-Go Configuration Validation Example ===")

	// Example 1: Load from JSON file
	fmt.Println("\n--- JSON Configuration Validation ---")

	// Create example config file
	exampleConfig := map[string]interface{}{
		"app": map[string]interface{}{
			"environment": "production",
			"debug":       false,
			"appName":     "MyAwesomeApp",
			"version":     "v1.2.3",
			"secretKey":   "supersecretkeythatisatleast32chars",
			"cors": map[string]interface{}{
				"allowedOrigins":   []string{"https://myapp.com", "https://admin.myapp.com"},
				"allowedMethods":   []string{"GET", "POST", "PUT", "DELETE"},
				"allowCredentials": true,
			},
			"rateLimit": map[string]interface{}{
				"enabled":        true,
				"requestsPerMin": 1000,
				"burst":          50,
			},
		},
		"server": map[string]interface{}{
			"host":        "0.0.0.0",
			"port":        8080,
			"readTimeout": 30,
			"tls": map[string]interface{}{
				"enabled":  true,
				"certFile": "/etc/ssl/certs/app.crt",
				"keyFile":  "/etc/ssl/private/app.key",
			},
		},
		"database": map[string]interface{}{
			"host":     "db.example.com",
			"port":     5432,
			"name":     "myapp_prod",
			"username": "dbuser",
			"password": "dbpassword",
			"sslMode":  "require",
		},
		"redis": map[string]interface{}{
			"host":     "redis.example.com",
			"port":     6379,
			"password": "redispassword",
			"database": 1,
		},
		"logging": map[string]interface{}{
			"level":  "info",
			"format": "json",
			"output": "file",
			"file": map[string]interface{}{
				"path":       "/var/log/myapp/app.log",
				"maxSize":    100,
				"maxBackups": 5,
				"compress":   true,
			},
		},
	}

	// Validate the example configuration
	if err := configSchema.Validate(exampleConfig); err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Example configuration is valid!")
		printConfigSummary(exampleConfig)
	}

	// Example 2: Test invalid configurations
	fmt.Println("\n--- Invalid Configuration Examples ---")

	invalidConfigs := []map[string]interface{}{
		{
			// Missing required fields
			"app": map[string]interface{}{
				"environment": "production",
				// Missing appName, version, secretKey
			},
			"server": map[string]interface{}{
				"port": 8080,
			},
			"database": map[string]interface{}{
				"host": "localhost",
				// Missing required database fields
			},
			"logging": map[string]interface{}{
				"level": "info",
			},
		},
		{
			"app": map[string]interface{}{
				"environment": "invalid_env", // Invalid environment
				"appName":     "MyApp",
				"version":     "1.0.0",
				"secretKey":   "short", // Too short
			},
			"server": map[string]interface{}{
				"port": 99999, // Invalid port
			},
			"database": map[string]interface{}{
				"host":     "localhost",
				"name":     "myapp",
				"username": "user",
				"password": "pass",
			},
			"logging": map[string]interface{}{
				"level": "trace", // Invalid log level
			},
		},
	}

	for i, config := range invalidConfigs {
		fmt.Printf("\nTesting invalid config #%d:\n", i+1)
		if err := configSchema.Validate(config); err != nil {
			if validationErr, ok := err.(*zod.ValidationError); ok {
				fmt.Printf("❌ %s\n", validationErr.ErrorJSON())
			} else {
				fmt.Printf("❌ %v\n", err)
			}
		} else {
			fmt.Println("⚠️  Unexpected: Invalid config passed validation!")
		}
	}

	// Example 3: Environment variables
	fmt.Println("\n--- Environment Variable Configuration ---")

	// Set some example environment variables
	os.Setenv("APP_NAME", "MyApp")
	os.Setenv("APP_VERSION", "v1.0.0")
	os.Setenv("APP_SECRET_KEY", "this_is_a_very_secure_secret_key_32")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "myapp")
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_PASSWORD", "password")

	if config, err := LoadConfigFromEnv(); err != nil {
		fmt.Printf("❌ Environment config validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Environment configuration is valid!")
		fmt.Printf("App: %s v%s\n", config.App.AppName, config.App.Version)
		fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
		fmt.Printf("Database: %s@%s:%d/%s\n",
			config.Database.Username, config.Database.Host,
			config.Database.Port, config.Database.Name)
	}
}

func printConfigSummary(config map[string]interface{}) {
	fmt.Println("\nConfiguration Summary:")

	if app, ok := config["app"].(map[string]interface{}); ok {
		fmt.Printf("  App: %s v%s (%s)\n",
			app["appName"], app["version"], app["environment"])
	}

	if server, ok := config["server"].(map[string]interface{}); ok {
		fmt.Printf("  Server: %s:%v\n", server["host"], server["port"])
	}

	if db, ok := config["database"].(map[string]interface{}); ok {
		fmt.Printf("  Database: %s@%s:%v/%s\n",
			db["username"], db["host"], db["port"], db["name"])
	}

	if redis, ok := config["redis"].(map[string]interface{}); ok {
		fmt.Printf("  Redis: %s:%v (db:%v)\n",
			redis["host"], redis["port"], redis["database"])
	}

	if logging, ok := config["logging"].(map[string]interface{}); ok {
		fmt.Printf("  Logging: %s level, %s format, output to %s\n",
			logging["level"], logging["format"], logging["output"])
	}
}
