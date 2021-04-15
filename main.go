package main

import (
	"api-simple-fiber/book"
	"api-simple-fiber/database"

	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//ValidateHeaders validate corp headers
func ValidateHeaders(c *fiber.Ctx) error {
	headersCopr := []string{"X-Txref", "X-Cmref", "X-Chref", "X-Country", "X-Commerce"}
	for _, header := range headersCopr {
		if err := c.Get(header); err == "" {
			log.Printf("Problema la cabecera %s", header)
			return c.Status(503).SendString("Service Unavailable")
		}
	}

	return c.Next()
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Use(ValidateHeaders)

	bookRoute := v1.Group("/book")
	bookRoute.Get("/", book.GetBooks)
	bookRoute.Post("/", book.NewBook)
	bookRoute.Get("/:id", book.GetBook)
	bookRoute.Patch("/:id", book.UpdateBook)
	bookRoute.Delete("/:id", book.DeleteBook)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "books.db")
	if err != nil {
		panic("Failed to connect database")
	}
	log.Println("Connection Opened to Database")
	database.DBConn.AutoMigrate(&book.Book{})
	log.Println("Database Migrated")
}

func main() {
	app := fiber.New()

	initDatabase()
	defer database.DBConn.Close()

	setupRoutes(app)
	log.Fatal(app.Listen(":3030"))
}
