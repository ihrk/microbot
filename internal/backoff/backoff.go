package backoff

import "time"

func RunWithRetry(retryLim int, retryBreak time.Duration, f func() error) error {
	var (
		err  error
		wait time.Duration
	)

	for i := 0; i < retryLim; i++ {
		err = f()
		if err == nil {
			break
		}

		time.Sleep(wait)

		wait = retryBreak << i
	}

	return err
}
