package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Peeranut-Kit/go_backend_test/handler"
	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/service"
	"github.com/Peeranut-Kit/go_backend_test/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	defer gracefulShutdown()
	fmt.Println("Hello")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	/*defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}*/
	fmt.Println("Database connected successfully")

	// AutoMigration to create task table in database. create but never delete column, so it is not practical. we preferred Migrator()
	db.AutoMigrate(&utils.Task{}, &utils.User{})

	// Initiate secondary adapter
	taskRepo := repo.NewTaskGormRepo(db)
	userRepo := repo.NewUserGormRepo(db)
	// Initiate primary adapter
	taskHandler := handler.NewHttpTaskHandler(taskRepo)
	userHandler := handler.NewHttpUserHandler(userRepo)
	
	// Fiber
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Enable CORS with default settings
	app.Use(cors.New())

	app.Use(simpleLogMiddleware)

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	// This way applies middleware into every request
	// Can group routes by using taskRoute := app.Group("/tasks")
	// taskRoute.Use(checkMiddleware) this only applies in one group
	// then taskRoute.Get("/", handler.GetTasks)

	taskRoute := app.Group("/tasks")
	taskRoute.Use(authRequiredMiddleware)

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

func initDatabase() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logger.Info, // Log level
		  IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
		  ParameterizedQueries:      false,           // Don't include params in the SQL log
		  Colorful:                  true,          // Disable color
		},
	)
	
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	//db, err := sql.Open("postgres", connStr)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
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

func simpleLogMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	return c.Next()
}

func authRequiredMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	secretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)
	fmt.Println(claim)

	return c.Next()
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
