package config

import "os"

const (
	defaultHost = "localhost"
	defaultPort = "8080"
)

type Config struct {
	// Settings of this server.
	ServerHost string
	ServerPort string

	// Settings of tracer agent (like datadog, jagger, etc).
	AgentHost string
	AgentPort string

	// Settings of database.
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// get configuration from environment variables.
func New() Config {
	var serverHost string
	if serverHost = os.Getenv("SERVER_HOST"); serverHost == "" {
		serverHost = defaultHost
	}

	var serverPort string
	if serverPort = os.Getenv("SERVER_PORT"); serverPort == "" {
		serverPort = defaultPort
	}

	agentHort := os.Getenv("AGENT_HOST")
	agentPort := os.Getenv("AGENT_PORT")

	var dbDriver string
	if dbDriver = os.Getenv("DB_DRIVER"); dbDriver == "" {
		dbDriver = "postgres"
	}

	var dbHost string
	if dbHost = os.Getenv("DB_HOST"); dbHost == "" {
		dbHost = "localhost"
	}

	var dbPort string
	if dbPort = os.Getenv("DB_PORT"); dbPort == "" {
		dbPort = "5432"
	}

	var dbUser string
	if dbUser = os.Getenv("DB_USER"); dbUser == "" {
		dbUser = "root"
	}

	var dbPassword string
	if dbPassword = os.Getenv("DB_PASSWORD"); dbPassword == "" {
		dbPassword = "root"
	}

	var dbName string
	if dbName = os.Getenv("DB_NAME"); dbName == "" {
		dbName = "postgres"
	}

	var dbSslMode string
	if dbSslMode = os.Getenv("DB_SSL_MODE"); dbSslMode == "" {
		dbSslMode = "disable"
	}

	return Config{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AgentHost:  agentHort,
		AgentPort:  agentPort,
		DBDriver:   dbDriver,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbSslMode,
	}
}
