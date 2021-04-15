package book

import (
	"api-simple-fiber/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Book struct {
	gorm.Model
	Title  string `json:"name"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(ctx *fiber.Ctx) error {
	db := database.DBConn
	var books []Book
	db.Find(&books)
	return ctx.JSON(books)
}

func GetBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.DBConn
	var book Book
	db.Find(&book, id)
	return ctx.JSON(book)
}

func NewBook(ctx *fiber.Ctx) error {
	db := database.DBConn
	book := new(Book)
	log.Println(book)
	if err := ctx.BodyParser(book); err != nil {
		return ctx.Status(503).SendString(err.Error())
	}
	db.Create(&book)
	return ctx.JSON(book)
}

func UpdateBook(ctx *fiber.Ctx) error {
	type DataUpdateBook struct {
		Title  string `json:"name"`
		Author string `json:"author"`
		Rating int    `json:"rating"`
	}

	var dataUB DataUpdateBook
	if err := ctx.BodyParser(&dataUB); err != nil {
		return ctx.Status(503).SendString(err.Error())
	}

	var book Book
	id := ctx.Params("id")
	db := database.DBConn
	db.First(&book, id)

	book = Book{
		Title: dataUB.Title,
		Author: dataUB.Author,
		Rating: dataUB.Rating,
	}

	db.Save(&book)
	return ctx.JSON(book)
}

func DeleteBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.DBConn

	var book Book
	db.First(&book, id)
	log.Println(&book)
	if book.Title == "" {
		return ctx.Status(404).SendString("Not Found the Book with ID")
	}

	db.Delete(&book)
	return ctx.Status(201).SendString("Book successfully deleted")
}
