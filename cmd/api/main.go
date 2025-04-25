package main

import (
	"log"
	"time"

	"garage-api/internal/config"
	"garage-api/internal/database"
	"garage-api/internal/handlers"
	"garage-api/internal/middleware"
	"garage-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "garage-api/docs"
)

// @title Garage API
// @version 1.0
// @description API for managing garage products and users
// @host 192.168.1.2
// @BasePath /api/v1
// @schemes https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	log.Println("🚀 Starting Garage API...")
	log.Println("📝 Loading environment variables...")
	
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: .env file not found: %v", err)
	} else {
		log.Println("✅ Environment variables loaded successfully")
	}

	// Load configuration
	log.Println("⚙️ Loading application configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}
	log.Printf("📊 Database Configuration:\n  Host: %s\n  Port: %d\n  User: %s\n  Database: %s\n  SSL Mode: %s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBSSLMode)

	// Initialize database
	log.Println("🔌 Connecting to database...")
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connection established successfully")

	// Initialize models
	productModel := &models.ProductModel{DB: db}
	productHandler := &handlers.ProductHandler{ProductModel: productModel}

	// Initialize router
	log.Println("🛠️ Setting up router...")
	router := gin.Default()

	// Add middleware
	log.Println("🔒 Adding middleware...")
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Swagger documentation
	log.Println("📚 Setting up Swagger documentation...")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("✅ Swagger documentation available at /swagger/index.html")

	// Public routes
	log.Println("🔓 Setting up public routes...")
	public := router.Group("/api/v1")
	{
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
		public.GET("/products", productHandler.GetAllProducts)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
		})
	}

	// Protected routes
	log.Println("🔐 Setting up protected routes...")
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/products", productHandler.CreateProduct)
		protected.GET("/products/:id", productHandler.GetProductByID)
		protected.PUT("/products/:id", productHandler.UpdateProduct)
		protected.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Start server
	port := "8080"
	log.Printf("🌐 Server starting on port %s...", port)
	log.Println("📡 API Endpoints:")
	log.Println("  🔓 Public:")
	log.Println("    POST /api/v1/register")
	log.Println("    POST /api/v1/login")
	log.Println("    GET  /api/v1/products")
	log.Println("  🔐 Protected:")
	log.Println("    POST   /api/v1/products")
	log.Println("    GET    /api/v1/products/:id")
	log.Println("    PUT    /api/v1/products/:id")
	log.Println("    DELETE /api/v1/products/:id")
	log.Println("  📚 Documentation:")
	log.Println("    GET /swagger/*any")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
} 