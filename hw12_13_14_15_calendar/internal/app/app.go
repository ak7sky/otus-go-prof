package app

import (
	"context"
	"time"
)

type App struct { // TODO
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, event *Event) (*Event, error)
	DeleteEvent(ctx context.Context, id string) error
	GetEventsForDay(ctx context.Context, date time.Time) ([]*Event, error)
	GetEventsForWeek(ctx context.Context, wStart time.Time) ([]*Event, error)
	GetEventsForMonth(ctx context.Context, mStart time.Time) ([]*Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
