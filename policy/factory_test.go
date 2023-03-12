package policy

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

func TestFactory(t *testing.T) {
	cfg := config.NewRateLimitConfig(1, time.Second)
	storage := keyvalue.NewInMemory()

	t.Run("fixed window", func(t *testing.T) {
		s, err := GetPolicy(FixedWindow, cfg, storage)
		require.NoError(t, err)
		expectedType := reflect.TypeOf(&fixedWindow{}).String()
		structType := reflect.TypeOf(s).String()

		require.Equal(t, expectedType, structType)
	})

	t.Run("token bucket", func(t *testing.T) {
		s, err := GetPolicy(TokenBucket, cfg, storage)
		require.NoError(t, err)
		expectedType := reflect.TypeOf(&tokenBucket{}).String()
		structType := reflect.TypeOf(s).String()

		require.Equal(t, expectedType, structType)
	})

	t.Run("unknown policy", func(t *testing.T) {
		_, err := GetPolicy(-1, cfg, storage)
		require.Error(t, err)
		require.ErrorIs(t, err, newUnknownPolicyError(-1))
		require.ErrorContains(t, err, "unknown rate limiting policy")
	})
}
