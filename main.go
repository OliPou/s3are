package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OliPou/s3are/internal/common"
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

	var ginRouterGroupName string = os.Getenv("GIN_ROUTER_GROUP_NAME")

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found in environment variables")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check if the database is ready
	if err := checkDatabase(db); err != nil {
		log.Fatal("Database is not ready:", err)
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

	v1Router := router.Group(fmt.Sprintf("/%s", ginRouterGroupName))
	v1Router.GET("/healthz", handlerHealthz)
	v1Router.POST("/uploadFileRequest", middleware.Auth(apiCfg.HandlerRequestUpload))
	v1Router.PUT("/fileUploaded", middleware.Auth(apiCfg.HandlerRequestUploadCompleted))
	v1Router.GET("/fileStatus", middleware.Auth(apiCfg.HandlerFileStatus))

	// Start the server
	if err := router.Run(":" + portString); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

// checkDatabase tries to ping the database until it succeeds or times out
func checkDatabase(db *sql.DB) error {
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		fmt.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database is not ready")
}

func handlerHealthz(c *gin.Context) {
	status := struct {
		Status string `json:"status"`
		Ready  string `json:"ready"`
	}{
		Status: "Ok",
		Ready:  "true",
	}
	common.RespondWithJSON(c, http.StatusOK, status)
}
