package ctx

import (
	"context"
	"time"
)

type ProvideWithCancel func() (context.Context, context.CancelFunc)

type ProvideWithTimeout func(timeout time.Duration) (context.Context, context.CancelFunc)

type ProvideWithDeadline func(deadline time.Time) (context.Context, context.CancelFunc)
