package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"n1h41/zolaris-backend-app/api/handlers"
	"n1h41/zolaris-backend-app/internal/aws"
	"n1h41/zolaris-backend-app/internal/config"
	"n1h41/zolaris-backend-app/internal/db"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/repositories"
	"n1h41/zolaris-backend-app/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize AWS clients
	log.Println("Initializing AWS clients...")
	awsClients, err := aws.InitAWSClients(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to initialize AWS clients: %v", err)
	}

	// Initialize database clients
	log.Println("Initializing database clients...")
	database, err := db.NewDatabase(awsClients.DynamoDB, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database clients: %v", err)
	}

	// Initialize repositories
	deviceRepo := repositories.NewDeviceRepository(database.GetDynamoClient())
	// Set the table names from configuration
	deviceRepo.WithTables(database.GetDeviceTableName(), database.GetMachineDataTableName())

	policyRepo := repositories.NewPolicyRepository(awsClients.GetIoTClient())

	// Initialize services
	deviceService := services.NewDeviceService(deviceRepo)
	policyService := services.NewPolicyService(policyRepo, cfg.AWS.IoTPolicy)

	// Initialize handlers
	addDeviceHandler := handlers.NewAddDeviceHandler(deviceService)
	attachIotPolicyHandler := handlers.NewAttachIotPolicyHandler(policyService)
	getDeviceSensorDataHandler := handlers.NewGetDeviceSensorDataHandler(deviceService)
	listUserDevicesHandler := handlers.NewListUserDevicesHandler(deviceService)

	// Create router with global middleware
	r := gin.New()

	// Apply global middleware
	r.Use(middleware.GinLoggerMiddleware())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Group private routes (require authentication)
	private := r.Group("/")
	private.Use(middleware.GinAuthMiddleware())
	{
		private.POST("/device/add", addDeviceHandler.HandleGin)
		private.GET("/user/devices", listUserDevicesHandler.HandleGin)
	}

	// Public routes (no authentication required)
	r.POST("/device/attach-policy", attachIotPolicyHandler.HandleGin)
	r.POST("/device/sensor-data", getDeviceSensorDataHandler.HandleGin)

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
