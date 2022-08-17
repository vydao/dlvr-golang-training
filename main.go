package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"desc"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	IsAvailable bool    `json:"is_available"`
}

var products []Product

func main() {
	products = []Product{
		{Name: "Product 1", Description: "Lorem ipsum dolor sit amet", Price: 1999, Rating: 4.5, IsAvailable: true},
		{Name: "Product 2", Description: "Lorem ipsum dolor sit amet", Price: 5000, Rating: 4.3, IsAvailable: true},
		{Name: "Product 3", Description: "Lorem ipsum dolor sit amet", Price: 289.5, Rating: 3.3, IsAvailable: true},
	}
	router := gin.New()
	groupV1 := router.Group("/api/v1")

	groupV1.Handle(http.MethodPost, "/products", func(c *gin.Context) {
		p := &Product{Name: "Product", Description: "Lorem ipsum", Price: 8383, Rating: 4.0, IsAvailable: true}
		products = append(products, *p)
		c.JSON(http.StatusOK, gin.H{"data": p})
	})

	groupV1.Handle(http.MethodGet, "/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": products})
	})

	groupV1.Handle(http.MethodDelete, "/products/:id", func(c *gin.Context) {
		// TODO
	})

	router.Run(":8080")
}
