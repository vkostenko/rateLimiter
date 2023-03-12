package policy

import (
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

type Policy int

const (
	FixedWindow Policy = iota
	TokenBucket
)

type RateLimitPolicy interface {
	IsAccepted(requestHash string) bool
}

func GetPolicy(
	policy Policy,
	config config.RateLimitConfig,
	storage keyvalue.Storage,
) (RateLimitPolicy, error) {
	switch policy {
	case FixedWindow:
		return newFixedWindow(config, storage), nil
	case TokenBucket:
		return newTokenBucket(config, storage), nil
	}

	return nil, newUnknownPolicyError(policy)
}
