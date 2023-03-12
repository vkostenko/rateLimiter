package rateLimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/policy"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

func TestNewRateLimiter(t *testing.T) {
	t.Run("success, request accepted", func(t *testing.T) {
		requestLimit := 3

		cfg := config.NewRateLimitConfig(requestLimit, time.Second)
		storage := keyvalue.NewInMemory()

		rateLimiter, err := NewRateLimiter(policy.FixedWindow, cfg, storage)
		require.NoError(t, err)

		for i := 0; i < requestLimit; i++ {
			result := rateLimiter.IsAccepted("customer.1")
			require.True(t, result)

			result = rateLimiter.IsAccepted("customer.2")
			require.True(t, result)
		}
	})

	t.Run("success, request blocked", func(t *testing.T) {
		cfg := config.NewRateLimitConfig(1, time.Second)
		storage := keyvalue.NewInMemory()

		rateLimiter, err := NewRateLimiter(policy.FixedWindow, cfg, storage)
		require.NoError(t, err)

		result := rateLimiter.IsAccepted("customer.1")
		require.True(t, result)

		// second attempt fails
		result = rateLimiter.IsAccepted("customer.1")
		require.False(t, result)
	})

	t.Run("fail: invalid policy", func(t *testing.T) {
		cfg := config.NewRateLimitConfig(1, time.Second)
		storage := keyvalue.NewInMemory()

		_, err := NewRateLimiter(-1, cfg, storage)
		require.Error(t, err)
		require.ErrorContains(t, err, "unknown")
		require.ErrorContains(t, err, "policy")
	})
}
