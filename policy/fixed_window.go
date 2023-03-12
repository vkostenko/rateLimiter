package policy

import (
	"time"

	"github.com/vkostenko/rateLimiter/config"
	"github.com/vkostenko/rateLimiter/storage/keyvalue"
)

type fixedWindow struct {
	config       config.RateLimitConfig
	stateStorage fixedWindowStateStorage
}

type fixedWindowState struct {
	CurrentRequestsCount int
	WindowStartTime      time.Time
}

type fixedWindowStateStorage interface {
	Exists(key string) bool
	Get(key string) fixedWindowState
	Set(key string, value fixedWindowState)
	Delete(key string)
}

func newFixedWindow(
	rateLimitConfig config.RateLimitConfig,
	storage keyvalue.Storage,
) *fixedWindow {
	adaptedStorage := keyvalue.NewStorageValueTypeDecorator[fixedWindowState](storage)

	return &fixedWindow{
		config:       rateLimitConfig,
		stateStorage: adaptedStorage,
	}
}

func (w *fixedWindow) IsAccepted(requestHash string) bool {
	state := w.getUpdatedState(requestHash)

	if state.CurrentRequestsCount >= w.config.GetLimit() {
		return false
	}

	state.CurrentRequestsCount++
	w.stateStorage.Set(requestHash, state)

	return true
}

func (w *fixedWindow) getUpdatedState(requestHash string) fixedWindowState {
	exists := w.stateStorage.Exists(requestHash)
	if !exists {
		return w.getEmptyWindow()
	}

	state := w.stateStorage.Get(requestHash)
	timeSinceWindowStart := time.Now().Sub(state.WindowStartTime)

	if timeSinceWindowStart >= w.config.GetInterval() {
		return w.getEmptyWindow()
	}

	return state
}

func (w *fixedWindow) getEmptyWindow() fixedWindowState {
	return fixedWindowState{
		CurrentRequestsCount: 0,
		WindowStartTime:      time.Now(),
	}
}
