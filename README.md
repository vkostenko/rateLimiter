Rate Limiter
====

Usage example:
```go
requestLimit := 3
interval := time.Second
cfg := config.NewRateLimitConfig(requestLimit, interval)
storage := keyvalue.NewInMemory()

rateLimiter, err := NewRateLimiter(policy.FixedWindow, cfg, storage)

result := rateLimiter.IsAccepted("customer.1.view_page")
```

Validation is done by unique request hash: requestHash

Configuration:
- Requests limit
- Interval

Policies available:
- Fixed Window
- Token Bucket

