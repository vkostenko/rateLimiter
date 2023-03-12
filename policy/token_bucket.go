package policy

import (
	"math"
	"time"

	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

type tokenBucket struct {
	config       config.RateLimitConfig
	stateStorage tokenBucketStateStorage
}

type tokenBucketState struct {
	CurrentTokensCount int
	LastRefillTime     time.Time
}

type tokenBucketStateStorage interface {
	Exists(key string) bool
	Get(key string) tokenBucketState
	Set(key string, value tokenBucketState)
	Delete(key string)
}

func newTokenBucket(
	rateLimitConfig config.RateLimitConfig,
	storage keyvalue.Storage,
) *tokenBucket {
	adaptedStorage := keyvalue.NewStorageValueTypeDecorator[tokenBucketState](storage)

	return &tokenBucket{
		config:       rateLimitConfig,
		stateStorage: adaptedStorage,
	}
}

func (b *tokenBucket) IsAccepted(requestHash string) bool {
	state := b.getRefilledBucket(requestHash)

	if state.CurrentTokensCount >= 1 {
		state.CurrentTokensCount--
		b.stateStorage.Set(requestHash, state)

		return true
	}

	return false
}

func (b *tokenBucket) getRefilledBucket(requestHash string) tokenBucketState {
	exists := b.stateStorage.Exists(requestHash)
	if !exists {
		return b.getFullBucket()
	}

	state := b.stateStorage.Get(requestHash)
	timeSinceRefill := time.Now().Sub(state.LastRefillTime)

	if timeSinceRefill >= b.config.GetInterval() {
		return b.getFullBucket()
	}

	requestLimit := float64(b.config.GetLimit())
	limitationInterval := float64(b.config.GetInterval())

	addTokens := math.Floor(requestLimit * float64(timeSinceRefill) / limitationInterval)
	if addTokens <= 0 {
		return state
	}

	timeToAdd := addTokens * limitationInterval / requestLimit
	refillTime := state.LastRefillTime.Add(time.Duration(timeToAdd))

	state.CurrentTokensCount += int(addTokens)
	state.LastRefillTime = refillTime

	return state
}

func (b *tokenBucket) getFullBucket() tokenBucketState {
	return tokenBucketState{
		CurrentTokensCount: b.config.GetLimit(),
		LastRefillTime:     time.Now(),
	}
}
