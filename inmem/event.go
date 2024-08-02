package inmem

import "sync"

type EventService struct {
	mux sync.RWMutex
}

func NewEventService() *EventService {
	return &EventService{}
}
