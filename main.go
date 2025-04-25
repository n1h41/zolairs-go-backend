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

	// Initialize database clients
	log.Println("Initializing database clients...")
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database clients: %v", err)
	}

	// Initialize repositories
	deviceRepo := repositories.NewDeviceRepository(database.GetDynamoClient())
	// Set the table names from configuration
	deviceRepo.WithTables(database.GetDeviceTableName(), database.GetMachineDataTableName())
	
	policyRepo := repositories.NewPolicyRepository(database.GetIoTClient())

	// Initialize services
	deviceService := services.NewDeviceService(deviceRepo)
	policyService := services.NewPolicyService(policyRepo, cfg.AWS.IoTPolicy)

	// Initialize handlers
	addDeviceHandler := handlers.NewAddDeviceHandler(deviceService)
	attachIotPolicyHandler := handlers.NewAttachIotPolicyHandler(policyService)
	getDeviceSensorDataHandler := handlers.NewGetDeviceSensorDataHandler(deviceService)
	listUserDevicesHandler := handlers.NewListUserDevicesHandler(deviceService)

	// Create router
	mux := http.NewServeMux()

	// Set up private routes (require authentication)
	mux.Handle("POST /device/add", applyMiddlewares(
		addDeviceHandler,
		middleware.AuthMiddleware,
	))

	mux.Handle("POST /device/sensor-data", applyMiddlewares(
		getDeviceSensorDataHandler,
	))

	mux.Handle("GET /user/devices", applyMiddlewares(
		listUserDevicesHandler,
		middleware.AuthMiddleware,
	))

	// Public routes (no authentication required)
	mux.Handle("POST /device/attach-policy", applyMiddlewares(
		attachIotPolicyHandler,
	))

	// Health check endpoint
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Create server with global middleware
	handler := applyMiddlewares(
		mux,
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
	)

	// Create server
	port := cfg.Server.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
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

// applyMiddlewares applies a chain of middleware to a handler
func applyMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

