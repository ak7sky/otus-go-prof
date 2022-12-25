package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	_ "github.com/jackc/pgx/stdlib" // need to use pgx driver
	"github.com/jmoiron/sqlx"
)

var (
	qCreateEvent = `INSERT INTO calendar.events (title, descr, event_start, event_end, notify_before)
					VALUES(:title, :descr, :start, :end, :notify)
					RETURNING id;`

	qSelectEvents = "SELECT * FROM calendar.events WHERE event_start >= :from AND event_start < :to"

	qUpdateEvent = `UPDATE calendar.events
					SET title = :title, descr = :descr, event_start = :start, event_end = :end, notify_before = :notify 
					WHERE id = :id`

	qDeleteEvent = "DELETE FROM calendar.events WHERE id = :id"
)

type rdbStorage struct {
	db *sqlx.DB
}

func newRdbStorage(ctx context.Context, dsn string) (*rdbStorage, error) {
	storage := &rdbStorage{}
	return storage, storage.connect(ctx, dsn)
}

func (storage *rdbStorage) CreateEvent(ctx context.Context, event *app.Event) (*app.Event, error) {
	rows, err := storage.db.NamedQueryContext(ctx, qCreateEvent, map[string]interface{}{
		"title":  event.Title,
		"descr":  event.Desc,
		"start":  event.Start,
		"end":    event.End,
		"notify": event.NotifyBefore,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event in db: %w", err)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get id of created event: %w", err)
	}

	return event, nil
}

func (storage *rdbStorage) UpdateEvent(ctx context.Context, event *app.Event) (*app.Event, error) {
	res, err := storage.db.NamedExecContext(ctx, qUpdateEvent, map[string]interface{}{
		"id":     event.ID,
		"title":  event.Title,
		"descr":  event.Desc,
		"start":  event.Start,
		"end":    event.End,
		"notify": event.NotifyBefore,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update event in db: %w", err)
	}

	updRows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to determine number of updated rows")
	}
	if updRows != 1 {
		return nil, fmt.Errorf("failed to update event in db, possible event '%s' doesn't exist", event.ID)
	}

	return event, nil
}

func (storage *rdbStorage) DeleteEvent(ctx context.Context, id string) error {
	res, err := storage.db.NamedExecContext(ctx, qDeleteEvent, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete event '%s' from db: %w", id, err)
	}
	delRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine number of deleted rows")
	}
	if delRows != 1 {
		return fmt.Errorf("failed to delete event from db, possible event '%s' doesn't exist", id)
	}

	return nil
}

func (storage *rdbStorage) GetEventsForDay(ctx context.Context, date time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, date, 24*time.Hour)
}

func (storage *rdbStorage) GetEventsForWeek(ctx context.Context, wStart time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, wStart, 7*24*time.Hour)
}

func (storage *rdbStorage) GetEventsForMonth(ctx context.Context, mStart time.Time) ([]*app.Event, error) {
	return storage.getEvents(ctx, mStart, 30*24*time.Hour)
}

func (storage *rdbStorage) getEvents(ctx context.Context, from time.Time, period time.Duration) ([]*app.Event, error) {
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	to := from.Add(period)

	rows, err := storage.db.NamedQueryContext(ctx, qSelectEvents, map[string]interface{}{"from": from, "to": to})
	if err != nil {
		return nil, fmt.Errorf("failed to select events from db: %w", err)
	}
	defer rows.Close()

	ee := make([]*app.Event, 0)
	for rows.Next() {
		var e app.Event
		err = rows.StructScan(&e)
		if err != nil {
			return nil, fmt.Errorf("failed to scan db data to app model: %w", err)
		}
		ee = append(ee, &e)
	}

	return ee, rows.Err()
}

func (storage *rdbStorage) connect(ctx context.Context, dsn string) (err error) {
	storage.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database by pgx driver : %w", err)
	}

	return storage.db.PingContext(ctx)
}

func (storage *rdbStorage) Close(ctx context.Context) error {
	return storage.db.Close()
}
