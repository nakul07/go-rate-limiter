package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitExceeded(t *testing.T) {
	config = LoadConfig()

	// Set up an HTTP request for the admin/1/dashboard endpoint
	req, err := http.NewRequest("GET", "/admin/1/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := RateLimitMiddleware(AdminDashboardHandler, "admin")

	// Perform 2 requests (the default limit for admin/1 is 2 )
	for i := 0; i < 3; i++ {
		// Create a ResponseRecorder to capture the HTTP response
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		println(rr.Code)
		println(i)
		if i < 2 && rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}

		if i >= 2 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status TooManyRequests, got %v", rr.Code)
		}

	}

}
