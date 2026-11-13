package main

import (
	"log"
	"net/http"
	"os"

	"libratrack/internal/app"
)

func main() {
	config := app.LoadConfig()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	log.Printf("✅ Connected to PostgreSQL: %s:%s/%s", config.DBHost, config.DBPort, config.DBName)
	log.Printf("📊 Connection pool: MaxOpenConns=%d, MaxIdleConns=%d", config.MaxOpenConn, config.MaxIdleConn)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("🚀 listening on %s", addr)
	if err := http.ListenAndServe(addr, app.NewRouterWithDB(db)); err != nil {
		log.Fatal(err)
	}
}
