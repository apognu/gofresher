package gofresher

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

type RefreshFunc[S any] func(*S) (*S, error)

type Gofresher[S any] struct {
	cacheDuration time.Duration

	flight      *singleflight.Group
	refreshFunc RefreshFunc[S]

	state       *S
	lastRefresh *time.Time
}

func NewGofresher[S any](cache time.Duration, refreshFunc RefreshFunc[S]) *Gofresher[S] {
	return &Gofresher[S]{
		flight:        &singleflight.Group{},
		cacheDuration: cache,
		refreshFunc:   refreshFunc,
		state:         nil,
	}
}

func (gr *Gofresher[S]) refresh() (*S, error) {
	if gr.refreshFunc == nil {
		return nil, fmt.Errorf("refresh function was not provided")
	}

	now := time.Now()

	newState, err, _ := gr.flight.Do("singleflight", func() (interface{}, error) {
		return gr.refreshFunc(gr.state)
	})

	if err != nil {
		return nil, err
	}

	if state, ok := newState.(*S); ok {
		gr.state = state
		gr.lastRefresh = &now

		return state, nil
	}

	return nil, fmt.Errorf("received invalid state from refresh function")
}

func (gr *Gofresher[S]) timedRefresh() (*S, error, bool) {
	now := time.Now()

	if gr.lastRefresh == nil || gr.lastRefresh.Add(gr.cacheDuration).Before(now) {
		state, err := gr.refresh()

		return state, err, true
	}

	return gr.state, nil, false
}

func (gr *Gofresher[S]) Start(tick time.Duration) {
	go func() {
		for {
			_, err, _ := gr.timedRefresh()

			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(tick)
		}
	}()
}

func (gr *Gofresher[S]) State(refreshable bool) (*S, error, bool) {
	if refreshable {
		return gr.timedRefresh()
	}

	return gr.state, nil, false
}

func (gr *Gofresher[S]) ForceRefresh() (*S, error) {
	return gr.refresh()
}
