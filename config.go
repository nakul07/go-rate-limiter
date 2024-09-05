package main

import (
	"sync"
	"time"
)

type EndpointConfig struct {
	MaxTokens  int           `json:"max_tokens"`
	RefillRate time.Duration `json:"refill_rate"`
}

type RateLimitConfig struct {
	mu          sync.Mutex
	RateLimiter *RateLimiter
	Metrics     *Metrics
	UserConfig  map[string]EndpointConfig
	AdminConfig map[string]EndpointConfig
}

// LoadConfig initializes the rate limit configuration.
func LoadConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RateLimiter: NewRateLimiter(),
		Metrics:     NewMetrics(),
		UserConfig: map[string]EndpointConfig{
			"default": {MaxTokens: 5, RefillRate: time.Minute / 5},
		},
		AdminConfig: map[string]EndpointConfig{
			"default": {MaxTokens: 2, RefillRate: time.Minute / 2},
		},
	}
}

// GetConfig updates the rate limit configuration for a given endpoint type and ID.
func (c *RateLimitConfig) UpdateConfig(endpointType, id string, maxTokens int, refillRate time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	config := EndpointConfig{MaxTokens: maxTokens, RefillRate: refillRate}

	if endpointType == "user" {
		c.UserConfig[id] = config
	} else if endpointType == "admin" {
		c.AdminConfig[id] = config
	}
}

// GetConfig retrieves the rate limit configuration for a given endpoint type and ID.
func (c *RateLimitConfig) GetConfig(endpointType, id string) EndpointConfig {
	c.mu.Lock()
	defer c.mu.Unlock()

	if endpointType == "user" {
		if config, exists := c.UserConfig[id]; exists {
			return config
		}
		return c.UserConfig["default"]
	} else if endpointType == "admin" {
		if config, exists := c.AdminConfig[id]; exists {
			return config
		}
		return c.AdminConfig["default"]
	}

	return EndpointConfig{}
}
