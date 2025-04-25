package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	transport "n1h41/zolaris-backend-app/internal/transport/http"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer recovery
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				
				// Return a 500 error to the client
				message := "Internal server error"
				if fmt.Sprintf("%v", err) != "" {
					message = fmt.Sprintf("Internal server error: %v", err)
				}
				
				transport.SendError(w, http.StatusInternalServerError, message)
			}
		}()
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

