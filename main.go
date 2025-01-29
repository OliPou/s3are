package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/OliPou/s3are/internal/database"
	"github.com/OliPou/s3are/middleware"
	s3uploadfile "github.com/OliPou/s3are/s3UploadFile"
	"github.com/OliPou/s3are/s3client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		// Continue execution as .env file might not exist in production
	}

	// Get PORT from environment variables with default fallback
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found in environment variables")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")
	if region == "" || bucket == "" {
		log.Fatal("AWS_REGION or S3_BUCKET not found in environment variables")
	}

	s3Client, err := s3client.NewS3Client(region, bucket)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := &s3uploadfile.ApiConfig{
		DB:       dbQueries,
		S3Client: s3Client,
	}

	fmt.Printf("Server starting on port: %s\n", portString)

	// Initialize the router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()

	// Get allowed origins from environment variable or use default
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// For development, you might want to allow all origins
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = []string{allowedOrigins}
	}

	// Additional CORS configurations
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * 60 * 60 // 12 hours

	// Add CORS middleware
	router.Use(cors.New(config))

	// Define routes
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	v1Router := router.Group("/v1")
	v1Router.GET("/healthz", s3uploadfile.HandlerHealthz)
	v1Router.POST("/uploadFileRequest", middleware.Auth(apiCfg.HandlerRequestUpload))
	v1Router.PUT("/fileUploaded", middleware.Auth(apiCfg.HandlerRequestUploadCompleted))
	v1Router.GET("/fileStatus", middleware.Auth(apiCfg.HandlerFileStatus))

	// Start the server
	if err := router.Run(":" + portString); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
