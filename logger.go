package main

import (
	"log"
	"net/http"
)

// LogRateLimitEvent logs a rate limit event.
func LogRateLimitEvent(r *http.Request) {
	log.Printf("Rate limit exceeded for %s", r.URL.Path)
}
