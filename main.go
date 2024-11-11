package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Peeranut-Kit/go_backend_test/handler"
	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	defer gracefulShutdown()
	fmt.Println("Hello")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected successfully")

	taskRepo := repo.NewPostgresDB(db)
	taskHandler := handler.TaskHandler{
		DB:       db,
		TaskRepo: taskRepo,
	}
	// Set up routes to handler
	http.HandleFunc("/tasks", taskHandler.TasksHandler)
	http.HandleFunc("/tasks/", taskHandler.TaskHandlerByID)

	// Start background task for periodic cleanup
	go service.BackgroundTask(taskRepo)

	// Start HTTP server
	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDatabase() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")
}