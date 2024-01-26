package gofresher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestManualRefresh(t *testing.T) {
	value := 42

	gr := NewGofresher[int](time.Minute, func(state *int) (*int, error) {
		value := value

		return &value, nil
	})

	state, err, refreshed := gr.State(true)

	assert.NotNil(t, state)
	assert.Nil(t, err)
	assert.True(t, refreshed)
	assert.Equal(t, value, *state)

	newState, err, refreshed := gr.State(true)

	assert.NotNil(t, newState)
	assert.Nil(t, err)
	assert.False(t, refreshed)
	assert.Equal(t, *state, *newState)
}

func TestForceRefresh(t *testing.T) {
	gr := NewGofresher[int](time.Minute, func(state *int) (*int, error) {
		value := 0

		if state != nil {
			value = *state + 1
		}

		return &value, nil
	})

	gr.ForceRefresh()
	gr.ForceRefresh()

	state, err := gr.ForceRefresh()

	assert.NotNil(t, state)
	assert.Nil(t, err)
	assert.Equal(t, 2, *state)
}

func TestConcurrentRefresh(t *testing.T) {
	updated := 0

	gr := NewGofresher[int](time.Minute, func(state *int) (*int, error) {
		time.Sleep(1 * time.Second)

		value := 42
		updated += 1

		return &value, nil
	})

	go gr.ForceRefresh()
	go gr.ForceRefresh()
	gr.ForceRefresh()

	state, err := gr.ForceRefresh()

	assert.Equal(t, 2, updated)
	assert.NotNil(t, state)
	assert.Nil(t, err)
	assert.Equal(t, 42, *state)
}
