package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI    string
	MongoDBName string
	JWTSecret   string
	HTTPPort    string
	FrontEndURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error, .env is optional

	cfg := &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		MongoDBName: os.Getenv("MONGO_DB_NAME"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		HTTPPort:    os.Getenv("HTTP_PORT"),
		FrontEndURL: os.Getenv("FRONTEND_URL"),
	}

	log.Default().Println("Configuration Loaded:" + cfg.MongoDBName)

	if cfg.MongoURI == "" {
		return nil, errors.New("MONGO_URI is required")
	}
	if cfg.MongoDBName == "" {
		return nil, errors.New("MONGO_DB_NAME is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}
	if cfg.FrontEndURL == "" {
		cfg.FrontEndURL = "http://localhost:5173"
	}

	return cfg, nil
}

// package config

// import (
// 	"errors"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type Config struct {
// 	MongoURI    string
// 	MongoDBName string
// 	JWTSecret   string
// 	HTTPPort    string
// }

// func Load() (*Config, error) {
// 	// Try current directory
// 	_ = godotenv.Load(".env")
// 	// Try one level above (when running from cmd/server)
// 	_ = godotenv.Load("../.env")

// 	cfg := &Config{
// 		MongoURI:    os.Getenv("MONGO_URI"),
// 		MongoDBName: os.Getenv("MONGO_DB_NAME"),
// 		JWTSecret:   os.Getenv("JWT_SECRET"),
// 		HTTPPort:    os.Getenv("HTTP_PORT"),
// 	}

// 	if cfg.MongoURI == "" {
// 		return nil, errors.New("MONGO_URI is required")
// 	}
// 	if cfg.MongoDBName == "" {
// 		return nil, errors.New("MONGO_DB_NAME is required")
// 	}
// 	if cfg.JWTSecret == "" {
// 		return nil, errors.New("JWT_SECRET is required")
// 	}
// 	if cfg.HTTPPort == "" {
// 		cfg.HTTPPort = "8080"
// 	}

// 	return cfg, nil
// }
