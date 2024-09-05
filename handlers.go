package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// UserDataHandler handles requests for user data.
func UserDataHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "User data accessed"}
	json.NewEncoder(w).Encode(response)
}

// AdminDashboardHandler handles requests for the admin dashboard.
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Admin dashboard accessed"}
	json.NewEncoder(w).Encode(response)
}

// PublicInfoHandler handles requests for public information.
func PublicInfoHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Public info accessed"}
	json.NewEncoder(w).Encode(response)
}

// UpdateRateLimitConfigHandler handles requests to update the rate limit configuration.
func UpdateRateLimitConfigHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EndpointType string `json:"endpoint_type"`
		ID           string `json:"id"`
		MaxTokens    int    `json:"max_tokens"`
		RefillRate   int    `json:"refill_rate_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.EndpointType != "user" && req.EndpointType != "admin" {
		http.Error(w, "Invalid endpoint type", http.StatusBadRequest)
		return
	}

	config.UpdateConfig(req.EndpointType, req.ID, req.MaxTokens, time.Duration(req.RefillRate)*time.Second)

	//Update the rate limit in Metrrics
	config.Metrics.SetRateLimit("/"+req.EndpointType+"/"+req.ID+"/dashboard", req.MaxTokens)

	// Reset the token bucket for this ID and endpoint type
	config.RateLimiter.ResetBucket(req.ID, req.EndpointType)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rate limit configuration updated successfully"))
}
