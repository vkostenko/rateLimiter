package config

import "time"

type RateLimitConfig struct {
	limit    int
	interval time.Duration
}

func NewRateLimitConfig(limit int, interval time.Duration) RateLimitConfig {
	return RateLimitConfig{
		limit:    limit,
		interval: interval,
	}
}

func (c *RateLimitConfig) GetLimit() int {
	return c.limit
}

func (c *RateLimitConfig) GetInterval() time.Duration {
	return c.interval
}
