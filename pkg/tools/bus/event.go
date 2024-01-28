package bus

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

var ErrStopPropagation = errors.New("propagation stopped")

type EventHandler[T any] func(e T) error

type eventStorage[T any] struct {
	id      string
	handler EventHandler[T]
}

type EventBag[T any] struct {
	mux      sync.RWMutex
	handlers []*eventStorage[T]
}

func (bus *EventBag[T]) Add(fn EventHandler[T]) string {
	bus.mux.Lock()
	defer bus.mux.Unlock()

	id := uuid.New().String()
	bus.handlers = append(bus.handlers, &eventStorage[T]{id, fn})

	return id
}

func (bus *EventBag[T]) Remove(id string) {
	bus.mux.Lock()
	defer bus.mux.Unlock()

	for i := len(bus.handlers) - 1; i >= 0; i-- {
		if bus.handlers[i].id == id {
			bus.handlers = append(bus.handlers[:i], bus.handlers[i+1:]...)
			return
		}
	}
}

func (bus *EventBag[T]) Trigger(e T, oneOffHandlers ...EventHandler[T]) error {
	bus.mux.RLock()

	handlers := make([]*eventStorage[T], 0, len(bus.handlers)+len(oneOffHandlers))
	handlers = append(handlers, bus.handlers...)

	// append the one-off handlers
	for i, oneOff := range oneOffHandlers {
		handlers = append(handlers, &eventStorage[T]{
			id:      fmt.Sprintf("@%d", i),
			handler: oneOff,
		})
	}

	bus.mux.RUnlock()

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
