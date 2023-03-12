package rateLimiter

import (
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/policy"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

type RateLimiter struct {
	config config.RateLimitConfig
	policy policy.RateLimitPolicy
}

func (r *RateLimiter) IsAccepted(requestHash string) bool {
	return r.policy.IsAccepted(requestHash)
}

func NewRateLimiter(
	policyName policy.Policy,
	config config.RateLimitConfig,
	storage keyvalue.Storage,
) (*RateLimiter, error) {
	rateLimiter, err := policy.GetPolicy(policyName, config, storage)
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		policy: rateLimiter,
		config: config,
	}, nil
}
