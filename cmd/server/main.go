package main

import (
	"billing/internal/api/routes"
	"billing/internal/db/ent"
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseCfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	SSLMode  string `json:"sslmode"`
}

func main() {
	client := initDatabase()
	defer closeDatabase(client)

	router := routes.SetupRoutes(client)
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func initDatabase() *ent.Client {
	cfg := &DatabaseCfg{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Type:     "postgres",
		SSLMode:  "disable",
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DbName, cfg.Password, cfg.SSLMode)

	client, err := ent.Open(cfg.Type, dataSourceName)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Enable debug logging
	client = client.Debug()

	return client
}

func closeDatabase(client *ent.Client) {
	if err := client.Close(); err != nil {
		log.Fatalf("failed closing connection to postgres: %v", err)
	}
}
