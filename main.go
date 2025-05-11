package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"n1h41/zolaris-backend-app/api/handlers"
	"n1h41/zolaris-backend-app/docs"
	"n1h41/zolaris-backend-app/internal/aws"
	"n1h41/zolaris-backend-app/internal/config"
	"n1h41/zolaris-backend-app/internal/db"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/repositories"
	"n1h41/zolaris-backend-app/internal/services"
)

func main() {
	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Swagger documentation
	docs.SwaggerInfo.Title = "Zolaris Backend API"
	docs.SwaggerInfo.Description = "API for IoT device management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize AWS clients
	log.Println("Initializing AWS clients...")
	awsClients, err := aws.InitAWSClients(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize AWS clients: %v", err)
	}

	// Initialize database clients
	log.Println("Initializing database clients...")
	database, err := db.NewDatabase(context.Background(), awsClients.DynamoDB, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database clients: %v", err)
	}

	// Initialize repositories
	deviceRepo := repositories.NewDeviceRepository(database.GetDynamoClient())
	policyRepo := repositories.NewPolicyRepository(awsClients.GetIoTClient())
	categoryRepo := repositories.NewCategoryRepository(database.GetDynamoClient())
	userRepo := repositories.NewUserRepository(database.GetDynamoClient())

	deviceRepo.WithTables(database.GetDeviceTableName(), database.GetMachineDataTableName())

	// Initialize services
	deviceService := services.NewDeviceService(deviceRepo)
	policyService := services.NewPolicyService(policyRepo, cfg.AWS.IoTPolicy)
	categoryService := services.NewCategoryService(categoryRepo)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	addDeviceHandler := handlers.NewAddDeviceHandler(deviceService)
	attachIotPolicyHandler := handlers.NewAttachIotPolicyHandler(policyService)
	getDeviceSensorDataHandler := handlers.NewGetDeviceSensorDataHandler(deviceService)
	listUserDevicesHandler := handlers.NewListUserDevicesHandler(deviceService)
	addCategoryHandler := handlers.NewAddCategoryHandler(categoryService)
	getCategoriesByTypeHandler := handlers.NewGetCategoriesByTypeHandler(categoryService)
	listAllCategoriesHandler := handlers.NewListAllCategoriesHandler(categoryService)
	checkParentIDHandler := handlers.NewCheckHasParentIDHandler(userService)
	updateUserDetailsHandler := handlers.NewUpdateUserDetailsHandler(userService)
	getUserDetailsHandler := handlers.NewGetUserDetailsHandler(userService)

	// Create router with global middleware
	r := gin.New()

	// Set up Swagger endpoint with dynamic host based on environment
	swaggerHost := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
	if cfg.Server.Environment == "production" || cfg.Server.Environment == "staging" {
		// In production/staging, use the actual host (could be set from config)
		swaggerHost = cfg.Server.ExternalURL
	}
	swaggerURL := ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", swaggerHost))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, swaggerURL))

	// Set up CORS
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow specific origins
			allowedOrigins := []string{
				"http://3.110.190.71",
				"https://staging.duvw6ii0xapud.amplifyapp.com",
			}
			if slices.Contains(allowedOrigins, origin) {
				return true
			}
			// Allow all localhost origins
			if len(origin) > 16 && origin[:16] == "http://localhost" {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	// Apply global middleware
	r.Use(middleware.GinLoggerMiddleware())
	r.Use(gin.Recovery())

	// Health check endpoint
	// @Summary Health check
	// @Description Check if the API is running
	// @Tags System
	// @Accept json
	// @Produce plain
	// @Success 200 {string} string "OK"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Group private routes (require authentication)
	private := r.Group("/")
	private.Use(middleware.GinAuthMiddleware())
	{
		private.POST("/device/add", addDeviceHandler.HandleGin)
		private.GET("/user/devices", listUserDevicesHandler.HandleGin)
		private.GET("/user/check-parent-id", checkParentIDHandler.HandleGin)

		// User details endpoints
		private.POST("/user/details", updateUserDetailsHandler.HandleGin)
		private.GET("/user/details", getUserDetailsHandler.HandleGin)
	}

	// Public routes (no authentication required)
	r.POST("/device/attach-policy", attachIotPolicyHandler.HandleGin)
	r.POST("/device/sensor-data", getDeviceSensorDataHandler.HandleGin)
	r.POST("/category/add", addCategoryHandler.HandleGin)
	r.GET("/category/type/:type", getCategoriesByTypeHandler.HandleGin)
	r.GET("/category/all", listAllCategoriesHandler.HandleGin)

	// Create server
	port := cfg.Server.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
