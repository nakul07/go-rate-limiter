# Rate Limiter HTTP Server

## Overview

This Go application implements a rate-limiting middleware using the Token Bucket algorithm. It supports different rate limits per endpoint and user, handles concurrency, and provides basic logging and metrics.

## Running the Service

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd go-rate-limiter
   ```

2. **Build the application:**

   ```bash
   go build -o go-rate-limiter
   ```

3. **Run the server:**

   ```bash
   ./go-rate-limiter
   ```

4. **Access the endpoints:**
   - `GET /user/:id/data`
   - `GET /admin/:id/dashboard`
   - `GET /public/info`
   - `GET /metrics`
   - `POST /config/update`
     - Example request body:
     ```json
     "endpoint_type": "admin",
     "id": "1",
     "max_tokens": 5,
     "refill_rate_seconds": 30
     ```

## Testing the System

You can run the tests with:

```bash
   go test ./
```
