package main

import (
	"log"
	"net/http"
	"time"
)

var config *RateLimitConfig

// Entry point of the application.
func main() {
	config = LoadConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/", RateLimitMiddleware(UserDataHandler, "user"))
	mux.HandleFunc("/admin/", RateLimitMiddleware(AdminDashboardHandler, "admin"))
	mux.HandleFunc("/public/info", PublicInfoHandler)
	mux.HandleFunc("/metrics", MetricsHandler)
	mux.HandleFunc("/config/update", UpdateRateLimitConfigHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}
