package storage

import (
	"context"
	"database/sql"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	mock              sqlmock.Sqlmock
	storage           rdbStorage
	qParamPlaceHolder = regexp.MustCompile(":(id|title|descr|start|end|notify|from|to)")
)

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setup() {
	var db *sql.DB
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Println("unexpected error while creating db mock")
	}
	dbx := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		log.Printf("unexpected error while init test data in db mock: %v", err)
	}

	storage = rdbStorage{db: dbx}
}

func TestRdbStorage_GetEventsForDay(t *testing.T) {
	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qSelectEvents, "?"))
	tm := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)

	expRows := sqlmock.NewRows([]string{"id", "title", "descr", "event_start", "event_end", "notify_before"}).
		AddRow("c35f1220-03d8-4931-85cb-c18665a55674", "ev1", "ev1 desc",
			tm, tm.Add(10*time.Minute), 900000000000).
		AddRow("19cc13a5-5d72-4fdb-b571-c0a2be2ef45d", "ev2", "ev2 desc",
			tm.Add(15*time.Minute), tm.Add(20*time.Minute), 800000000000)

	expFrom := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	expTo := expFrom.Add(24 * time.Hour)
	mock.ExpectQuery(expQuery).WithArgs(expFrom, expTo).WillReturnRows(expRows)

	ee, err := storage.GetEventsForDay(context.Background(), tm)

	require.NoError(t, err)
	require.Equal(t, 2, len(ee))
}

func TestRdbStorage_GetEventsForWeek(t *testing.T) {
	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qSelectEvents, "?"))
	tm := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)

	expRows := sqlmock.NewRows([]string{"id", "title", "descr", "event_start", "event_end", "notify_before"}).
		AddRow("c35f1220-03d8-4931-85cb-c18665a55674", "ev1", "ev1 desc",
			tm, tm.Add(10*time.Minute), 900000000000).
		AddRow("3ba11bd3-4aad-4c71-8dea-45e80da3ddb9", "ev3", "ev3 desc",
			tm.Add(3*24*time.Hour), tm.Add(3*24*time.Hour+15*time.Minute), 900000000000)

	expFrom := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	expTo := expFrom.Add(7 * 24 * time.Hour)
	mock.ExpectQuery(expQuery).WithArgs(expFrom, expTo).WillReturnRows(expRows)

	ee, err := storage.GetEventsForWeek(context.Background(), tm)

	require.NoError(t, err)
	require.Equal(t, 2, len(ee))
}

func TestRdbStorage_GetEventsForMonth(t *testing.T) {
	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qSelectEvents, "?"))
	tm := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)

	expRows := sqlmock.NewRows([]string{"id", "title", "descr", "event_start", "event_end", "notify_before"}).
		AddRow("c35f1220-03d8-4931-85cb-c18665a55674", "ev1", "ev1 desc",
			tm, tm.Add(10*time.Minute), 900000000000).
		AddRow("3ba11bd3-4aad-4c71-8dea-45e80da3ddb9", "ev3", "ev3 desc",
			tm.Add(15*24*time.Hour), tm.Add(15*24*time.Hour+15*time.Minute), 900000000000)

	expFrom := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	expTo := expFrom.Add(30 * 24 * time.Hour)
	mock.ExpectQuery(expQuery).WithArgs(expFrom, expTo).WillReturnRows(expRows)

	ee, err := storage.GetEventsForMonth(context.Background(), tm)

	require.NoError(t, err)
	require.Equal(t, 2, len(ee))
}

func TestRdbStorage_CreateEvent(t *testing.T) {
	ev := &app.Event{
		Title:        "event to create",
		Desc:         "event description",
		Start:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		End:          time.Date(2023, 1, 1, 0, 10, 0, 0, time.Local),
		NotifyBefore: 15 * time.Minute,
	}

	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qCreateEvent, "?"))

	expCreatedID := "c35f1220-03d8-4931-85cb-c18665a55674"
	expRows := sqlmock.NewRows([]string{"id"}).AddRow(expCreatedID)
	mock.ExpectQuery(expQuery).WithArgs(ev.Title, ev.Desc, ev.Start, ev.End, ev.NotifyBefore).WillReturnRows(expRows)

	require.Equal(t, "", ev.ID)
	ev, err := storage.CreateEvent(context.Background(), ev)
	require.NoError(t, err)
	require.Equal(t, expCreatedID, ev.ID)
}

func TestRdbStorage_UpdateEvent(t *testing.T) {
	inEv := &app.Event{
		ID:           "c35f1220-03d8-4931-85cb-c18665a55674",
		Title:        "upd event",
		Desc:         "event description",
		Start:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		End:          time.Date(2023, 1, 1, 0, 10, 0, 0, time.Local),
		NotifyBefore: 15 * time.Minute,
	}

	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qUpdateEvent, "?"))

	// Success
	mock.ExpectExec(expQuery).
		WithArgs(inEv.Title, inEv.Desc, inEv.Start, inEv.End, inEv.NotifyBefore, inEv.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	outEv, err := storage.UpdateEvent(context.Background(), inEv)
	require.NoError(t, err)
	require.Equal(t, inEv, outEv)

	// Fail
	mock.ExpectExec(expQuery).
		WithArgs(inEv.Title, inEv.Desc, inEv.Start, inEv.End, inEv.NotifyBefore, inEv.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	outEv, err = storage.UpdateEvent(context.Background(), inEv)
	require.Error(t, err)
	require.Nil(t, outEv)
}

func TestRdbStorage_DeleteEvent(t *testing.T) {
	id := "c35f1220-03d8-4931-85cb-c18665a55674"
	expQuery := regexp.QuoteMeta(qParamPlaceHolder.ReplaceAllLiteralString(qDeleteEvent, "?"))

	// Success
	mock.ExpectExec(expQuery).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
	err := storage.DeleteEvent(context.Background(), id)
	require.NoError(t, err)

	// Fail
	mock.ExpectExec(expQuery).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
	err = storage.DeleteEvent(context.Background(), id)
	require.Error(t, err)
}
