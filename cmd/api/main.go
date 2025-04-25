package main

import (
	"fmt"
	"log"
	"net/http"
	"garage-api/docs"
	"garage-api/internal/database"
	"garage-api/internal/handlers"
	"garage-api/internal/middleware"
	"garage-api/internal/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Garage API
// @version         1.0
// @description     A product management API with authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@garage.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body handlers.loginRequest true "User registration"
// @Success 201 {object} handlers.loginResponse
// @Router /auth/register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Login user
// @Description Login user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body handlers.loginRequest true "User login"
// @Success 200 {object} handlers.loginResponse
// @Router /auth/login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func getAllProductsHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Get a product
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func getProductHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.Product true "Product to create"
// @Success 201 {object} models.Product
// @Router /products [post]
func createProductHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product to update"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func updateProductHandler(w http.ResponseWriter, r *http.Request) {}

// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 204
// @Router /products/{id} [delete]
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	// Database configuration
	dbConfig := &database.Config{
		Host:     "192.168.1.2",
		Port:     5432,
		User:     "casaos",
		Password: "casaos",
		DBName:   "garage-web",
	}

	// Connect to database
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	// Initialize models
	productModel := &models.ProductModel{DB: db}
	userModel := &models.UserModel{DB: db}

	// Initialize handlers
	productHandler := &handlers.ProductHandler{Model: productModel}
	authHandler := &handlers.AuthHandler{UserModel: userModel}

	// Create router
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	apiRouter.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	apiRouter.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Public product routes
	apiRouter.HandleFunc("/products", productHandler.GetAll).Methods("GET")

	// Protected product routes (require JWT)
	productRouter := apiRouter.PathPrefix("/products").Subrouter()
	productRouter.Use(middleware.AuthMiddleware)

	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.Get).Methods("GET")
	productRouter.HandleFunc("", productHandler.Create).Methods("POST")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.Update).Methods("PUT")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.Delete).Methods("DELETE")

	// Swagger documentation
	docs.SwaggerInfo.Title = "Garage API"
	docs.SwaggerInfo.Description = "A product management API with authentication"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
} 