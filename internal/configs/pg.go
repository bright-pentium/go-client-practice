package configs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Query    map[string]string
}

// LoadPostgresConfig reads from .env and builds a config struct.
func LoadPostgresConfig(envFilePath string) (*PostgresConfig, error) {
	envMap, err := godotenv.Read(envFilePath)
	if err != nil {
		return nil, err
	}

	rawPort := getEnv(envMap, "POSTGRES_PORT", "5432")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, err
	}

	query := make(map[string]string)
	queryJSON := getEnv(envMap, "POSTGRES_QUERY", "{}")
	if err := json.Unmarshal([]byte(queryJSON), &query); err != nil {
		return nil, err
	}

	return &PostgresConfig{
		Host:     getEnv(envMap, "POSTGRES_HOST", "mavis-postgres"),
		Port:     port,
		Database: getEnv(envMap, "POSTGRES_DATABASE", "mavis"),
		Username: getEnv(envMap, "POSTGRES_USER", ""),
		Password: getEnv(envMap, "POSTGRES_PASSWORD", ""),
		Query:    query,
	}, nil
}

func (c *PostgresConfig) DBURL() string {
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(c.Username, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Database,
	}

	// Add query parameters
	q := url.Values{}
	for key, val := range c.Query {
		q.Set(key, val)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// Helpers

func getEnv(env map[string]string, key, defaultVal string) string {
	if val, ok := env[key]; ok {
		return val
	}
	return defaultVal
}
