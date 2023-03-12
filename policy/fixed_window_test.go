package policy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

func TestFixedWindow_RateLimit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		requestLimit := 3
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second)
		storage := keyvalue.NewInMemory()
		validator := newFixedWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)

			result = validator.IsAccepted("customer.2")
			require.True(t, result, "iteration: %d", i)
		}
	})

	t.Run("success with full refill", func(t *testing.T) {
		requestLimit := 3
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second/10)
		storage := keyvalue.NewInMemory()
		validator := newFixedWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)

			result = validator.IsAccepted("customer.2")
			require.True(t, result, "iteration: %d", i)
		}

		time.Sleep(time.Second / 10)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)

			result = validator.IsAccepted("customer.2")
			require.True(t, result, "iteration: %d", i)
		}
	})

	t.Run("fail for partial refill", func(t *testing.T) {
		requestLimit := 5
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second/10)
		timeForOneBucket := time.Second / (10 * 5)

		storage := keyvalue.NewInMemory()
		validator := newFixedWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)
		}

		time.Sleep(timeForOneBucket)

		result := validator.IsAccepted("customer.1")
		require.False(t, result)
	})

	t.Run("fail", func(t *testing.T) {
		requestLimit := 3
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second)
		storage := keyvalue.NewInMemory()
		validator := newFixedWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result)
		}

		result := validator.IsAccepted("customer.1")
		require.False(t, result)
	})
}
