package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type EndpointMetrics struct {
	TotalRequests       int `json:"total_requests"`
	RateLimitedRequests int `json:"rate_limited_requests"`
	RateLimit           int `json:"rate_limit"`
}

type Metrics struct {
	mu      sync.Mutex
	metrics map[string]*EndpointMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: make(map[string]*EndpointMetrics),
	}
}

// IncrementTotal increments the total requests for the given endpoint.
func (m *Metrics) IncrementTotal(endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.metrics[endpoint]; !exists {
		m.metrics[endpoint] = &EndpointMetrics{}
	}
	m.metrics[endpoint].TotalRequests++
}

// IncrementRateLimited increments the rate-limited requests for the given endpoint.
func (m *Metrics) IncrementRateLimited(endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.metrics[endpoint]; !exists {
		m.metrics[endpoint] = &EndpointMetrics{}
	}
	m.metrics[endpoint].RateLimitedRequests++
}

// SetRateLimit sets the rate limit for the given endpoint.
func (m *Metrics) SetRateLimit(endpoint string, rateLimit int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.metrics[endpoint]; !exists {
		m.metrics[endpoint] = &EndpointMetrics{}
	}
	m.metrics[endpoint].RateLimit = rateLimit
}

// GetMetrics returns the metrics for all endpoints.
func (m *Metrics) GetMetrics() map[string]*EndpointMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.metrics
}

// MetricsHandler handles the /metrics endpoint.
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := config.Metrics.GetMetrics()
	json.NewEncoder(w).Encode(metrics)
}
