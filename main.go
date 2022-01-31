package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type purse struct {
	amount float64
}

var purse1 = purse{
	amount: 100,
}

var books = []book{
	{ID: 1, Title: "Book 1", Author: "Author 1", Price: 15.2, Quantity: 1},
	{ID: 2, Title: "Book 2", Author: "Author 2", Price: 25.2, Quantity: 2},
	{ID: 3, Title: "Book 3", Author: "Author 3", Price: 35.2, Quantity: 3},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var newBook book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id int) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	//convert string to int
	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not available"})
		return
	}

	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	book.Quantity++
	c.IndentedJSON(http.StatusOK, book)
}

func buyBookById(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not available"})
		return
	}

	purse1.amount += book.Price
	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func sellBookById(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if book.Price*0.75 > purse1.amount {
		c.JSON(http.StatusConflict, gin.H{"error": "Not enough money"})
		return
	}

	purse1.amount -= book.Price * 0.75
	book.Quantity++
	c.IndentedJSON(http.StatusOK, book)
}

func removeBookById(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	//convert string to int
	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := getBookById(bookId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	max_number := len(books)
	// save the book before deleting it
	book_to_delete := *book

	index := 0
	for i, b := range books {
		if b.ID == book.ID {
			books = append(books[:i], books[i+1:]...)
			break
		}
		index++
	}
	if index == max_number {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book_to_delete)
}

func getPurse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, purse1.amount)
}

func main() {
	router := gin.Default()
	router.GET("/balance", getPurse)
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", addBook)
	router.POST("/books/buy", buyBookById)
	router.POST("/books/sell", sellBookById)
	router.PATCH("/books/checkout", checkoutBook)
	router.PATCH("/books/return", returnBook)
	router.DELETE("/books/burn", removeBookById)
	router.Run(":8080")
}
