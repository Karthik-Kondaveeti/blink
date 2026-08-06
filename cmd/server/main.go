package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Karthik-Kondaveeti/blink/internal/cache/redis"
	"github.com/Karthik-Kondaveeti/blink/internal/database/postgres"
	"github.com/Karthik-Kondaveeti/blink/internal/generator/snowflake"
	"github.com/Karthik-Kondaveeti/blink/internal/handler"
	"github.com/Karthik-Kondaveeti/blink/internal/logger"
	"github.com/Karthik-Kondaveeti/blink/internal/service"
	"github.com/Karthik-Kondaveeti/blink/migration"
)

func main() {
	logger := logger.New()

	connStr := os.Getenv("DATABASE_URL")
	databaseTableName := os.Getenv("TABLE_NAME")
	db, err := postgres.New(connStr, databaseTableName)
	if err != nil {
		log.Fatal("error creating new database", err)
		return
	}
	defer db.Close()

	if err := migration.RunMigrations(connStr); err != nil {
		log.Fatal("error running migrations", err)
	}

	workerID, err := strconv.ParseUint(os.Getenv("WORKER_ID"), 10, 16)
	if err != nil {
		log.Fatal("error converting WORKER_ID to integer", err)
		return
	}

	generator, err := snowflake.New(uint16(workerID))
	if err != nil {
		log.Fatal("error creating new generator", err)
		return
	}

	cacheURL := os.Getenv("CACHE_URL")
	cachePASS := os.Getenv("CACHE_PASS")
	cache, err := redis.New(cacheURL, cachePASS, 0)
	if err != nil {
		log.Fatal("error creating new cache database", err)
		return
	}
	defer cache.Close()

	service, err := service.New(db, cache, generator, logger)
	if err != nil {
		log.Fatal("error creating new service", err)
		return
	}

	handler, err := handler.New(service)
	if err != nil {
		log.Fatal("error creating new server", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten/", handler.AddLinkHandler)
	mux.HandleFunc("GET /{shortCode}", handler.GetLinkHandler)

	Port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Fatal("error converting Port to integer", err)
		return
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%v", Port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	log.Printf("Server started on port %v", Port)
}
