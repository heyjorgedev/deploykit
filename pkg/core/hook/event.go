package hook

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

var ErrStopPropagation = errors.New("propagation stopped")

type Handler[T any] func(e T) error

type handlerPair[T any] struct {
	id      string
	handler Handler[T]
}

type Hook[T any] struct {
	mux      sync.RWMutex
	handlers []*handlerPair[T]
}

func (h *Hook[T]) Add(fn Handler[T]) string {
	h.mux.Lock()
	defer h.mux.Unlock()

	id := uuid.New().String()
	h.handlers = append(h.handlers, &handlerPair[T]{id, fn})

	return id
}

func (h *Hook[T]) Remove(id string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for i := len(h.handlers) - 1; i >= 0; i-- {
		if h.handlers[i].id == id {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

func (h *Hook[T]) Trigger(e T, oneOffHandlers ...Handler[T]) error {
	h.mux.RLock()

	handlers := make([]*handlerPair[T], 0, len(h.handlers)+len(oneOffHandlers))
	handlers = append(handlers, h.handlers...)

	// append the one-off handlers
	for i, oneOff := range oneOffHandlers {
		handlers = append(handlers, &handlerPair[T]{
			id:      fmt.Sprintf("@%d", i),
			handler: oneOff,
		})
	}

	h.mux.RUnlock()

	for _, h := range handlers {
		if err := h.handler(e); err != nil {
			if errors.Is(err, ErrStopPropagation) {
				return nil
			}
			return err
		}
	}

	return nil
}
