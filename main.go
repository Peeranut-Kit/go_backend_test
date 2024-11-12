package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Peeranut-Kit/go_backend_test/handler"
	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/service"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
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
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected successfully")

	taskRepo := repo.NewPostgresDB(db)
	taskHandler := handler.TaskHandler{
		TaskRepo: taskRepo,
	}

	// Fiber
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Enable CORS with default settings
	app.Use(cors.New())

	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))
	app.Use(checkMiddleware)

	// This way applies middleware into every request
	// Can group routes by using taskRoute := app.Group("/tasks")
	// taskRoute.Use(checkMiddleware) this only applies in one group
	// then taskRoute.Get("/", handler.GetTasks)

	app.Get("/tasks", taskHandler.GetTasksHandler)
	app.Post("/tasks", taskHandler.PostTaskHandler)
	app.Get("/tasks/:id", taskHandler.GetTaskHandler)
	app.Put("/tasks/:id", taskHandler.PutTaskHandler)
	app.Delete("/tasks/:id", taskHandler.DeleteTaskHandler)

	// View Template -> render webpage without using frontend framework (no more usage)
	app.Get("/view-tasks", func(c *fiber.Ctx) error {
		return c.Render("task-index", fiber.Map{
			"Title":   "Task List",
			"Content": "[task1, task2, ...]",
		})
	})

	app.Get("/config", getEnv)

	// Start HTTP server
	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port %s...\n", port)
	app.Listen(":" + port)

	// Start background task for periodic cleanup
	go service.BackgroundTask(taskRepo)
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

func checkMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["admin"] != true {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

type User struct {
	Email    string
	Password string
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	var user *User
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !(user.Email == memberUser.Email && user.Password == memberUser.Password) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  user.Email,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   t,
	})
}

func getEnv(c *fiber.Ctx) error {
	// os.LookupEnv() looks for env in local machine
	if value, exist := os.LookupEnv("SECRET"); exist {
		return c.JSON(fiber.Map{
			"SECRET": value,
		})
	}
	return c.JSON(fiber.Map{
		"SECRET": "defaultSecret",
	})
}
