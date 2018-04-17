package health

import (
	"time"

	"github.com/pkg/errors"
)

const (
	defaultRecheckLimit    = 120
	defaultRecheckInterval = 1 * time.Second
)

// NewWatcher creates a new watcher with useful defaults
func NewWatcher() *Watcher {
	return &Watcher{
		RecheckLimit:    defaultRecheckLimit,
		RecheckInterval: defaultRecheckInterval,
		ResultListener: func(counter int, err error) {
			// do nothing
		},
	}
}

// CheckResultListener is always called with the latest result from the executor. The CheckResultListener can be used
// to display the results to the user or for logging purposes.
type CheckResultListener func(counter int, err error)

// Watcher is able to block until a checker becomes healthy
type Watcher struct {
	RecheckLimit    int
	RecheckInterval time.Duration
	ResultListener  CheckResultListener
}

// WaitUntilHealthy blocks the routine until checker returns no error. The watcher will check the result of the checker
// in the configured time interval until the limit is reached or the checker returns nil.
func (watcher *Watcher) WaitUntilHealthy(checker HealthChecker) error {
	ticker := time.NewTicker(watcher.RecheckInterval)
	defer ticker.Stop()

	counter := 0

	for range ticker.C {

		err := checker()
		if err == nil {
			return nil
		}

		counter++

		// call listener on error
		watcher.ResultListener(counter, err)

		if counter >= watcher.RecheckLimit {
			return errors.New("limit reached durring wait checker to become healthy")
		}
	}

	return errors.New("loop break out. This part should never be reached.")
}
