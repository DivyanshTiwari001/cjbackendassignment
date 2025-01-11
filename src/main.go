package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/DivyanshTiwari001/cjbackend/src/db"
	"github.com/DivyanshTiwari001/cjbackend/src/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	db.ConnectDB()

	// disconnect 
	defer db.CloseDB()

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
        AllowOrigins: os.Getenv("CORS_ORIGIN"), 
        AllowMethods: "GET,POST,PUT,DELETE",                         
        AllowHeaders: "Origin, Content-Type, Accept",              
        AllowCredentials: true,           
    }))

	// Routes
	routes.SetupRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(app.Listen(":" + port))
}
