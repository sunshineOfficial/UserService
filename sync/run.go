package sync

import "context"

type WaitFn func() error

func WaitContext(ctx context.Context, fn WaitFn) error {
	c := make(chan error, 1)

	go func(c chan error, fn func() error) {
		c <- fn()
		close(c)
	}(c, fn)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
