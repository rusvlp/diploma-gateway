package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}

	GatewayPort    string
	UserGRPCAddr   string
	CourseGRPCAddr string
	JWTSecret      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file didn't found")
	}

	cfg := &Config{}

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Database = os.Getenv("POSTGRES_DB")

	cfg.GatewayPort = os.Getenv("GATEWAY_PORT")
	cfg.UserGRPCAddr = os.Getenv("USER_GRPC_ADDR")
	cfg.CourseGRPCAddr = os.Getenv("COURSE_GRPC_ADDR")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	return cfg
}
