package main

import (
	"os"
	"github.com/akifanabil/synapsis-backend-challenge/handlers"
	// "github.com/akifanabil/synapsis-backend-challenge/middleware"
	"github.com/akifanabil/synapsis-backend-challenge/migrations"
	"github.com/akifanabil/synapsis-backend-challenge/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"	
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Title Backend Engineer Challenge API Documentation
// Version 1.0
// Description API Documentation for Backend Engineer Challenge
// ContactName Akifa Nabil Ufairah

// Server localhost:8080 Server 1
// Server localhost:8081 Server 2

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}

	HOST := os.Getenv("API_HOST")
	PORT := os.Getenv("API_PORT")

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Horizon APP API Documentation"
	docs.SwaggerInfo.Description = "API Documentation for Horizon App"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
	
	DB := migrations.Init()
	h := handlers.New(DB)

	router := gin.Default()

	api := router.Group(docs.SwaggerInfo.BasePath)
	api.POST("/login", h.Login)
	api.POST("/register", h.Register)
	api.GET("/product/:category", h.GetProducts)
	api.GET("/cart",h.GetCart)
	api.POST("/cart", h.AddCart)
	api.DELETE("/cart", h.DeleteCartItem)
	api.POST("/checkout", h.Checkout)
	
	

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	router.Run(HOST+":"+PORT)
}
