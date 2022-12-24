package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetPostgresConnectionStr() string {
	config := Config{
		Host:     GoDotEnvVariable("DB_HOST"),
		Port:     GoDotEnvVariable("DB_PORT"),
		User:     GoDotEnvVariable("DB_USER"),
		Password: GoDotEnvVariable("DB_PASS"),
		DBName:   GoDotEnvVariable("DB_NAME"),
		SSLMode:  GoDotEnvVariable("DB_SSLMODE"),
	}

	postgresConnectionStr :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	return postgresConnectionStr
}
