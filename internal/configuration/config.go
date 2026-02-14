package configuration

import (
	"log/slog"
	"os"

	"github.com/RodrigoMattosoSilveira/WorkIncomeExpenses/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	APP_ENV         string
	DB_NAME         string
	PROXY_PORT      string
	GIN_PORT        string
	FIBER_PORT      string
	CSRF_SECRET     string
	SESSION_KEY     string
	JWT_KEY         string
	TMPL_ROOT       string
	PERSON_ROLES    string
	PERSON_STATUSES string
}

var Cfg *Config

func LoadConfig() error {

	homeDir, err := utils.FindProjectRoot()
	if err != nil {
		slog.Error("Error calculating project's home directory")
	}

	envFile := homeDir + "/" + ".env"
	err = godotenv.Load(envFile)
	if err != nil {
		slog.Error("Error loading .env file")
		panic(err)
	}

	envSecretsFile := homeDir + "/" + ".env.secrets"
	err = godotenv.Load(envSecretsFile)
	if err != nil {
		slog.Error("Error loading .env.secrets file")
		panic(err)
	}

	Cfg = &Config{
		//                    key           default value, if key not found
		APP_ENV:         GetEnv("APP_ENV", "development"),
		DB_NAME:         GetEnv("DB_NAME", "/private/var/ContasCorrentes/sqlite_dev.db"),
		FIBER_PORT:      GetEnv("FIBER_PORT", "3000"),
		CSRF_SECRET:     GetEnv("CSRF_SECRET", "default-secret-must-be-32-chars-long"),
		SESSION_KEY:     GetEnv("SESSION_KEY", "default-secret-must-be-32-chars-long"),
		JWT_KEY:         GetEnv("JWT_KEY", "default-secret-must-be-32-chars-long"),
		TMPL_ROOT:       GetEnv("TMPL_ROOT", "internal/templates/root_new"),
		PERSON_ROLES:    GetEnv("PERSON_ROLES", "System,Tenant,Application,Operator,Person"),
		PERSON_STATUSES: GetEnv("PERSON_STATUSES", "Active,Inactive"),
	}
	slog.Info("Configuration loaded", "app_env", Cfg.APP_ENV)
	return nil
}

func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
