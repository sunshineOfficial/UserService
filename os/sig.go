package os

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WaitTerminate(mainCtx context.Context, quitFn func(ctx context.Context)) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if quitFn == nil {
		return
	}

	quitCtx, cancelQuitCtx := context.WithTimeout(mainCtx, 15*time.Second)
	defer cancelQuitCtx()

	quitFn(quitCtx)
}
