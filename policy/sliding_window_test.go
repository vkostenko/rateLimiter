package policy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

func TestSlidingWindow_RateLimit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		requestLimit := 3
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second)
		storage := keyvalue.NewInMemory()
		validator := newSlidingWindow(rlConfig, storage)

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
		validator := newSlidingWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)

			result = validator.IsAccepted("customer.2")
			require.True(t, result, "iteration: %d", i)
		}

		time.Sleep(2 * time.Second / 10)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)

			result = validator.IsAccepted("customer.2")
			require.True(t, result, "iteration: %d", i)
		}
	})

	t.Run("partial refill", func(t *testing.T) {
		requestLimit := 5
		interval := time.Second / 10
		rlConfig := config.NewRateLimitConfig(requestLimit, interval)
		timeForOneBucket := interval / time.Duration(requestLimit)

		storage := keyvalue.NewInMemory()
		validator := newSlidingWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result, "iteration: %d", i)
		}

		time.Sleep(timeForOneBucket)

		// same time window, hence fails
		result := validator.IsAccepted("customer.1")
		require.False(t, result)

		// new time window + enough time passed for 1 request
		time.Sleep(time.Duration(requestLimit) * timeForOneBucket)

		result = validator.IsAccepted("customer.1")
		require.True(t, result)

		result = validator.IsAccepted("customer.1")
		require.False(t, result)
	})

	t.Run("fail", func(t *testing.T) {
		requestLimit := 3
		rlConfig := config.NewRateLimitConfig(requestLimit, time.Second)
		storage := keyvalue.NewInMemory()
		validator := newSlidingWindow(rlConfig, storage)

		for i := 0; i < requestLimit; i++ {
			result := validator.IsAccepted("customer.1")
			require.True(t, result)
		}

		result := validator.IsAccepted("customer.1")
		require.False(t, result)
	})
}
