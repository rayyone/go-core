package retry

import (
	"time"

	"github.com/rayyone/go-core/errors"
	loghelper "github.com/rayyone/go-core/helpers/log"
)

type Options struct {
	DelayBetweenAttempt time.Duration
	MaxRetry            int
	attempt             int
}

func DefaultOptions() Options {
	return Options{
		DelayBetweenAttempt: 1 * time.Second,
		MaxRetry:            3,
	}
}

func WithRetry(fn func() error, opts Options) error {
	if opts.attempt > 0 {
		loghelper.PrintYellowf(
			"Retrying %d/%d in %.2f s..", opts.attempt, opts.MaxRetry, opts.DelayBetweenAttempt.Seconds(),
		)
		time.Sleep(opts.DelayBetweenAttempt)
	}
	err := fn()
	if err != nil {
		opts.attempt++
		if opts.attempt > opts.MaxRetry {
			if opts.MaxRetry > 0 {
				return errors.BadRequest.Newf("Max retry reached. Error: %+v", err)
			} else {
				return errors.BadRequest.Newf(err.Error())
			}
		} else {
			return WithRetry(fn, opts)
		}
	}

	return nil
}
