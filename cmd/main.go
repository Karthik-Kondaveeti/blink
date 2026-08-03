package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Karthik-Kondaveeti/blink/internal/cache/redis"
	snowflake "github.com/Karthik-Kondaveeti/blink/internal/generator/snowflake"
	"github.com/Karthik-Kondaveeti/blink/internal/handler"
	"github.com/Karthik-Kondaveeti/blink/internal/logger"
	"github.com/Karthik-Kondaveeti/blink/internal/service"
	"github.com/Karthik-Kondaveeti/blink/internal/storage/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env", err)
		return
	}

	logger := logger.New()

	databaseAddress := os.Getenv("DB_URL")
	databaseTableName := os.Getenv("TABLE_NAME")
	database, err := postgres.New(databaseAddress, databaseTableName)
	if err != nil {
		log.Fatal("error creating new database", err)
		return
	}
	defer database.Close()

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

	cacheAddress := os.Getenv("CACHE_URL")
	cachePassword := os.Getenv("CACHE_PASS")
	cache, err := redis.New(cacheAddress, cachePassword, 0)
	if err != nil {
		log.Fatal("error creating new cache database", err)
		return
	}
	defer cache.Close()

	service, err := service.New(database, cache, generator, logger)
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

	Port, err := strconv.Atoi(os.Getenv("Port"))
	if err != nil {
		log.Fatal("error converting Port to integer", err)
		return
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%v", Port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	log.Printf("Server started on port %v", Port)
}
