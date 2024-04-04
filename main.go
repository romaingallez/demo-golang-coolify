package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("demo-golang-coolify: ")
	log.SetOutput(os.Stderr)

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// add additional code here
}

func main() {

	err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	DB.AutoMigrate(&Book{}, &Author{})

	InitDemoDB()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	e.GET("/books", GetBooks)

	e.POST("/books", CreateBook)

	e.GET("/books/:id", GetBook)

	e.GET("/authors", GetAuthors)

	e.POST("/authors", CreateAuthor)

	e.GET("/authors/:id", GetAuthor)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))

}

// Handler
func hello(c echo.Context) error {
	return c.String(200, "Hello, World!")
}

func GetBooks(c echo.Context) error {
	var books []Book
	DB.Preload("Authors").Find(&books)
	return c.JSON(200, books)
}

func CreateBook(c echo.Context) error {
	book := Book{}
	c.Bind(&book)
	DB.Create(&book)
	return c.JSON(200, book)
}

func GetBook(c echo.Context) error {
	id := c.Param("id")
	var book Book
	DB.First(&book, id)
	return c.JSON(200, book)
}

func CreateAuthor(c echo.Context) error {
	author := Author{}
	c.Bind(&author)
	DB.Create(&author)
	return c.JSON(200, author)
}

func GetAuthors(c echo.Context) error {
	var authors []Author
	DB.Preload("Books").Find(&authors)
	return c.JSON(200, authors)
}

func GetAuthor(c echo.Context) error {
	id := c.Param("id")
	var author Author
	DB.First(&author, id)
	return c.JSON(200, author)
}

func ConnectDB() (err error) {
	uri := os.Getenv("DB_URI")
	if uri == "" {
		return fmt.Errorf("DB_URI is not set")

	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris", host, username, password, strings.TrimPrefix(u.Path, "/"), port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func InitDemoDB() {
	// Test if both Author and Book tables are empty
	var authors []Author
	DB.Find(&authors)
	var books []Book
	DB.Find(&books)

	if len(authors) == 0 && len(books) == 0 {
		// Create Authors and their Books
		Authors := []Author{
			{
				Name: "Philip K. Dick",
				Books: []Book{
					{Title: "Do Androids Dream of Electric Sheep?", Length: 256, Language: "English"},
					{Title: "The Man in the High Castle", Length: 324, Language: "English"},
				},
			},
			{
				Name: "George Orwell",
				Books: []Book{
					{Title: "1984", Length: 328, Language: "English"},
					{Title: "Animal Farm", Length: 112, Language: "English"},
				},
			},
		}
		DB.Create(Authors)

	}

}
