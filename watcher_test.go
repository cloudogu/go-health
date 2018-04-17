package health_test

import (
	"testing"

	"time"

	"github.com/cloudogu/go-health"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWatcher_WaitUntilHealthy_Success(t *testing.T) {
	watcher := health.NewWatcher()

	err := watcher.WaitUntilHealthy(func() error {
		return nil
	})
	assert.NoError(t, err)
}

func TestWatcher_WaitUntilHealthy_BecomesHealthy(t *testing.T) {
	watcher := health.NewWatcher()
	watcher.RecheckLimit = 10
	watcher.RecheckInterval = 10 * time.Millisecond

	counter := 0
	err := watcher.WaitUntilHealthy(func() error {
		if counter < 5 {
			counter++
			return errors.New("error")
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestWatcher_WaitUntilHealthy_LimitReached(t *testing.T) {
	lastSeenCounter := -1

	watcher := health.NewWatcher()
	watcher.RecheckLimit = 3
	watcher.RecheckInterval = 10 * time.Millisecond
	watcher.ResultListener = func(counter int, err error) {
		lastSeenCounter = counter
	}

	err := watcher.WaitUntilHealthy(func() error {
		return errors.New("error")
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "limit")
	assert.Equal(t, 3, lastSeenCounter)
}
