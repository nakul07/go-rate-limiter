package main

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
}

type RateLimiter struct {
	mu           sync.Mutex
	userBuckets  map[string]*TokenBucket
	adminBuckets map[string]*TokenBucket
}

// NewTokenBucket creates a new token bucket with the given maximum tokens and refill rate.
func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// AllowRequest checks if a request is allowed based on the token bucket's state.
func (tb *TokenBucket) AllowRequest() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tb.tokens += int(elapsed / tb.refillRate)
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		userBuckets:  make(map[string]*TokenBucket),
		adminBuckets: make(map[string]*TokenBucket),
	}
}

// GetBucket returns the token bucket for the given ID and bucket type.
func (rl *RateLimiter) GetBucket(id string, bucketType string, maxTokens int, refillRate time.Duration) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	var bucket *TokenBucket
	if bucketType == "user" {
		if _, exists := rl.userBuckets[id]; !exists {
			rl.userBuckets[id] = NewTokenBucket(maxTokens, refillRate)
		}
		bucket = rl.userBuckets[id]
	} else if bucketType == "admin" {
		if _, exists := rl.adminBuckets[id]; !exists {
			rl.adminBuckets[id] = NewTokenBucket(maxTokens, refillRate)
		}
		bucket = rl.adminBuckets[id]
	}

	return bucket
}

// ResetBucket resets the token bucket for the given ID and bucket type.
func (rl *RateLimiter) ResetBucket(id string, bucketType string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Get the updated config
	config := rl.getRateLimitConfig(bucketType, id)

	// Reset the token bucket with new settings
	if bucketType == "user" {
		rl.userBuckets[id] = NewTokenBucket(config.MaxTokens, config.RefillRate)
	} else if bucketType == "admin" {
		rl.adminBuckets[id] = NewTokenBucket(config.MaxTokens, config.RefillRate)
	}
}

// Helper method to get the current rate limit config
func (rl *RateLimiter) getRateLimitConfig(bucketType string, id string) EndpointConfig {
	if bucketType == "user" {
		return config.GetConfig("user", id)
	} else if bucketType == "admin" {
		return config.GetConfig("admin", id)
	}
	return EndpointConfig{}
}
