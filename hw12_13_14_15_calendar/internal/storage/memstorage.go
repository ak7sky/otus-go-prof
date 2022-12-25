package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/google/uuid"
)

var (
	errBlankTitle   = errors.New("event title must be filled")
	errBlankStart   = errors.New("event start time must be filled")
	errBlankEnd     = errors.New("event end time must be filled")
	errReservedTime = errors.New("time is reserved by another event")
	errNotFound     = errors.New("event not found")
)

type memStorage struct {
	mu         sync.RWMutex
	eventsByID map[string]*app.Event
}

func newMemStorage(ctx context.Context) (*memStorage, error) {
	return &memStorage{
		mu:         sync.RWMutex{},
		eventsByID: make(map[string]*app.Event),
	}, nil
}

func (storage *memStorage) CreateEvent(ctx context.Context, event *app.Event) (*app.Event, error) {
	if err := storage.validateEvent(event); err != nil {
		return nil, err
	}

	var id string
	for {
		id = uuid.NewString()
		if !storage.isEventFound(id) {
			break
		}
	}

	event.ID = id
	storage.mu.Lock()
	storage.eventsByID[id] = event
	storage.mu.Unlock()
	return event, nil
}

func (storage *memStorage) UpdateEvent(ctx context.Context, event *app.Event) (*app.Event, error) {
	if !storage.isEventFound(event.ID) {
		return nil, errNotFound
	}

	if err := storage.validateEvent(event); err != nil {
		return nil, err
	}

	storage.mu.Lock()
	storage.eventsByID[event.ID] = event
	storage.mu.Unlock()

	return event, nil
}

func (storage *memStorage) DeleteEvent(ctx context.Context, id string) error {
	if !storage.isEventFound(id) {
		return errNotFound
	}

	storage.mu.Lock()
	delete(storage.eventsByID, id)
	storage.mu.Unlock()

	return nil
}

func (storage *memStorage) GetEventsForDay(ctx context.Context, date time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, date, 24*time.Hour)
}

func (storage *memStorage) GetEventsForWeek(ctx context.Context, wStart time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, wStart, 7*24*time.Hour)
}

func (storage *memStorage) GetEventsForMonth(ctx context.Context, mStart time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, mStart, 30*24*time.Hour)
}

func (storage *memStorage) getEvents(ctx context.Context, from time.Time, period time.Duration) ([]*app.Event, error) {
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	to := from.Add(period)
	ee := make([]*app.Event, 0)

	storage.mu.RLock()
	for _, e := range storage.eventsByID {
		if e.Start.Equal(from) || (e.Start.After(from) && e.Start.Before(to)) {
			ee = append(ee, e)
		}
	}
	storage.mu.RUnlock()

	return ee, nil
}

func (storage *memStorage) validateEvent(event *app.Event) error {
	if event.Title == "" {
		return errBlankTitle
	}
	if event.Start.IsZero() {
		return errBlankStart
	}
	if event.End.IsZero() {
		return errBlankEnd
	}

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	for _, existing := range storage.eventsByID {
		if event.ID == existing.ID {
			continue
		}
		if event.Start.Equal(existing.Start) ||
			event.Start.Equal(existing.End) ||
			event.End.Equal(existing.Start) ||
			event.End.Equal(existing.End) ||
			(event.Start.After(existing.Start) && event.Start.Before(existing.End)) ||
			(event.End.After(existing.Start) && event.End.Before(existing.End)) ||
			(event.Start.Before(existing.Start) && event.End.After(existing.End)) {
			return errReservedTime
		}
	}

	return nil
}

func (storage *memStorage) isEventFound(id string) bool {
	storage.mu.RLock()
	_, found := storage.eventsByID[id]
	storage.mu.RUnlock()
	return found
}
