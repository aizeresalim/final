package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tenajuro12/newBackend/db"
	"github.com/tenajuro12/newBackend/routes"
)

func main() {
	db.Connect()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	port := os.Getenv("PORT")
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":" + port)

}
