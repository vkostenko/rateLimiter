package policy

import (
	"math"
	"time"

	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

type slidingWindow struct {
	config       config.RateLimitConfig
	stateStorage slidingWindowStateStorage
}

type slidingWindowState struct {
	PreviousRequestsCount int
	CurrentRequestsCount  int
	WindowStartTime       time.Time
}

type slidingWindowStateStorage interface {
	Exists(key string) bool
	Get(key string) slidingWindowState
	Set(key string, value slidingWindowState)
	Delete(key string)
}

func newSlidingWindow(
	rateLimitConfig config.RateLimitConfig,
	storage keyvalue.Storage,
) *slidingWindow {
	adaptedStorage := keyvalue.NewStorageValueTypeDecorator[slidingWindowState](storage)

	return &slidingWindow{
		config:       rateLimitConfig,
		stateStorage: adaptedStorage,
	}
}

func (w *slidingWindow) IsAccepted(requestHash string) bool {
	state := w.getUpdatedState(requestHash)

	timeSinceWindowStart := time.Now().Sub(state.WindowStartTime)
	currentWindowRatio := float64(timeSinceWindowStart) / float64(w.config.GetInterval())

	requestsInPreviousWindow := (1 - currentWindowRatio) * float64(state.PreviousRequestsCount)
	requestsSpent := state.CurrentRequestsCount + int(math.Ceil(requestsInPreviousWindow))

	if requestsSpent >= w.config.GetLimit() {
		return false
	}

	state.CurrentRequestsCount++
	w.stateStorage.Set(requestHash, state)

	return true
}

func (w *slidingWindow) getUpdatedState(requestHash string) slidingWindowState {
	exists := w.stateStorage.Exists(requestHash)
	if !exists {
		return w.getEmptyWindow()
	}

	state := w.stateStorage.Get(requestHash)
	timeSinceWindowStart := time.Now().Sub(state.WindowStartTime)

	if timeSinceWindowStart >= 2*w.config.GetInterval() {
		return w.getEmptyWindow()
	}

	if timeSinceWindowStart >= w.config.GetInterval() {
		return slidingWindowState{
			PreviousRequestsCount: state.CurrentRequestsCount,
			CurrentRequestsCount:  0,
			WindowStartTime:       state.WindowStartTime.Add(w.config.GetInterval()),
		}
	}

	return state
}

func (w *slidingWindow) getEmptyWindow() slidingWindowState {
	return slidingWindowState{
		PreviousRequestsCount: 0,
		CurrentRequestsCount:  0,
		WindowStartTime:       time.Now(),
	}
}
