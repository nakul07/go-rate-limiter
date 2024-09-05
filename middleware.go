package main

import (
	"net/http"
	"strings"
)

// RateLimitMiddleware is a middleware function that applies rate limiting to incoming requests.
func RateLimitMiddleware(next http.HandlerFunc, bucketType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		id := parts[2]
		endpoint := r.URL.Path

		// Track the total requests for the endpoint
		config.Metrics.IncrementTotal(endpoint)

		// Get the configuration for this endpoint and ID
		endpointConfig := config.GetConfig(bucketType, id)

		// Set rate limit of the endpoint
		config.Metrics.SetRateLimit(endpoint, endpointConfig.MaxTokens)

		bucket := config.RateLimiter.GetBucket(id, bucketType, endpointConfig.MaxTokens, endpointConfig.RefillRate)

		// Apply rate limiting
		if bucket.AllowRequest() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			LogRateLimitEvent(r)
			config.Metrics.IncrementRateLimited(endpoint)
		}
	}
}
