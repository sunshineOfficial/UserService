package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const (
	listenTimeout   = 3 * time.Second
	shutdownTimeout = 5 * time.Second
)

type Server interface {
	Start()
	Stop()
	UseHandler(http.Handler)
}

type HTTPServer struct {
	ctx     context.Context
	log     *zap.Logger
	server  *http.Server
	running *atomic.Bool
}

func NewHTTPServer(ctx context.Context, log *zap.Logger, addr string) *HTTPServer {
	return &HTTPServer{
		ctx: ctx,
		log: log,
		server: &http.Server{
			Addr: addr,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		running: &atomic.Bool{},
	}
}

func (h *HTTPServer) Start() {
	if h.running.Load() {
		return
	}

	h.running.Store(true)
	go h.listen()
}

func (h *HTTPServer) listen() {
	h.log.Debug("Server is starting")

	for h.running.Load() {
		err := h.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Debug(fmt.Sprintf("Failed to start server. Retry in %s", listenTimeout), zap.Error(err))
			continue
		}
	}

	h.log.Debug("Server is stopped")
}

func (h *HTTPServer) Stop() {
	if !h.running.Load() {
		return
	}

	h.running.Store(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := h.server.Shutdown(shutdownCtx); err != nil {
		h.log.Debug(fmt.Sprintf("Failed to shutdown server in %s", shutdownTimeout), zap.Error(err))
	}
}

func (h *HTTPServer) UseHandler(handler http.Handler) {
	h.server.Handler = handler
}
