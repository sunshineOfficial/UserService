package kafka

import (
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

type listener struct {
	log         *zap.Logger
	consume     func() (Message, error)
	subs        []Subscriber
	subsMutex   *sync.Mutex
	subsRunning *atomic.Bool
}

func newListener(log *zap.Logger, consume func() (Message, error)) *listener {
	return &listener{
		log:         log,
		consume:     consume,
		subsMutex:   &sync.Mutex{},
		subsRunning: &atomic.Bool{},
	}
}

func (l *listener) start() {
	if l.subsRunning.Load() {
		return
	}

	l.subsRunning.Store(true)

	go l.listen()
}

func (l *listener) listen() {
	for l.subsRunning.Load() {
		msg, err := l.consume()
		l.broadcastMessage(msg, err)
	}
}

func (l *listener) broadcastMessage(message Message, err error) {
	l.subsMutex.Lock()
	defer l.subsMutex.Unlock()

	for _, s := range l.subs {
		l.sendMessage(s, message, err)
	}
}

func (l *listener) sendMessage(s Subscriber, message Message, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			l.log.Error("could not process message", zap.String("key", string(message.Key)), zap.Error(err), zap.Any("panic", panicErr))
		}
	}()

	s(message, err)
}

func (l *listener) stop() {
	l.subsRunning.Store(false)
}

func (l *listener) add(s Subscriber) {
	l.subsMutex.Lock()
	defer l.subsMutex.Unlock()

	l.subs = append(l.subs, s)
}
