package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ovh-prep/models"
	// "ovh-prep/workflows"
	// "go.temporal.io/sdk/client"
	// "go.temporal.io/sdk/worker"
)

var DB *gorm.DB
// var temporalClient client.Client
// var temporalWorker worker.Worker

func connectDB() {
	cnnString := "host=localhost user=postgres password=postgres dbname=bookstore port=5433"
	var err error
	DB, err = gorm.Open(postgres.Open(cnnString), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect", err)
	}

	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Failed to migrate", err)
	}

	log.Println("Connected to database")
}

// func setupTemporal() {
// 	var err error

// 	temporalClient, err = client.Dial(client.Options{
// 		HostPort: "localhost:7233", 
// 	})
// 	if err != nil {
// 		log.Fatalf("Unable to create Temporal client: %v", err)
// 	}

// 	temporalWorker = worker.New(temporalClient, "book-reservation-task-queue", worker.Options{})

// 	temporalWorker.RegisterWorkflow(workflows.ReserveBookWorkflow)
// 	temporalWorker.RegisterActivity(workflows.CheckBookAvailability)
// 	temporalWorker.RegisterActivity(workflows.ReserveBook)

// 	go func() {
// 		err := temporalWorker.Start()
// 		if err != nil {
// 			log.Fatalf("Unable to start Temporal worker: %v", err)
// 		}
// 	}()
// }

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Create(&book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func GetBooks(c *gin.Context) {
	var books []models.Book
	DB.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func GetBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	if err := DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func UpdateBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	if err := DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Save(&book)
	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated",
		"data":    book,
	})
}

func ReserveBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	var requestData struct {
		Reserver string `json:"reserver"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println("JSON Binding Error:", err) // Ajoute ce log pour voir l'erreur exacte
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.Quantity == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not available"})
		return
	}

	if book.IsReserved {
		if book.Reserver == requestData.Reserver {
			c.JSON(http.StatusConflict, gin.H{"message": "Book is already reserved by you"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"message": "Book is already reserved"})
		}
		return
	}

	book.IsReserved = true
	book.Reserver = requestData.Reserver
	book.Quantity -= 1

	DB.Save(&book)
	c.JSON(http.StatusOK, gin.H{
		"message": "Book reserved",
		"data":    book,
	})
}

func DeleteBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	if err := DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func SetupRouter(router *gin.Engine) {
	router.POST("/books", CreateBook)
	router.GET("/books", GetBooks)
	router.GET("/books/:id", GetBook)
	router.PUT("/books/:id", UpdateBook)
	router.DELETE("/books/:id", DeleteBook)
	router.POST("/books/reserve/:id", ReserveBook)
}

func main() {
	connectDB()

	router := gin.Default()
	SetupRouter(router)

	// setupTemporal()

	router.Run(":8089")
}
