package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimitConfig(t *testing.T) {
	requestLimit := 10
	interval := 5 * time.Second

	config := NewRateLimitConfig(requestLimit, interval)

	require.Equal(t, requestLimit, config.GetLimit())
	require.Equal(t, interval, config.GetInterval())
}
