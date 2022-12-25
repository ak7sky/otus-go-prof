package storage

import (
	"context"
	"testing"
	"time"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_CreateEvent(t *testing.T) {
	storage, err := newMemStorage(context.Background())
	require.NoError(t, err, "unexpected error during mem storage init")

	exEvID := "da632c1e-ffc5-463d-8b8f-6cdfbd70053d"
	exEvStart := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)
	exEvEnd := exEvStart.Add(15 * time.Minute)

	storage.eventsByID[exEvID] = &app.Event{ID: exEvID, Title: "existing ev", Start: exEvStart, End: exEvEnd}

	testCases := []struct {
		name        string
		event       *app.Event
		expectedErr error
	}{
		{
			name:        "blank title err",
			event:       &app.Event{Title: ""},
			expectedErr: errBlankTitle,
		},
		{
			name:        "blank start time",
			event:       &app.Event{Title: "new ev"},
			expectedErr: errBlankStart,
		},
		{
			name:        "blank end time",
			event:       &app.Event{Title: "new ev", Start: time.Now().Add(15 * time.Minute)},
			expectedErr: errBlankEnd,
		},
		{
			name:        "reserved time err",
			event:       &app.Event{Title: "new ev", Start: exEvStart, End: exEvEnd},
			expectedErr: errReservedTime,
		},
		{
			name:  "successful create",
			event: &app.Event{Title: "new ev", Start: exEvEnd.Add(1 * time.Minute), End: exEvEnd.Add(10 * time.Minute)},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, 1, len(storage.eventsByID))
			_, err := storage.CreateEvent(context.Background(), tc.event)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, 1, len(storage.eventsByID))
			} else {
				require.NoError(t, err, "unexpected error")
				require.Equal(t, 2, len(storage.eventsByID))
			}
		})
	}
}

func TestMemStorage_GetEvents(t *testing.T) {
	storage, err := newMemStorage(context.Background())
	require.NoError(t, err, "unexpected error during mem storage init")

	storage.eventsByID["9fae9b58-a03f-446d-8274-e121b8960b80"] = &app.Event{
		ID:    "9fae9b58-a03f-446d-8274-e121b8960b80",
		Title: "ev1",
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		End:   time.Date(2023, 1, 1, 0, 30, 0, 0, time.Local),
	}
	storage.eventsByID["bf268c55-11e0-4efc-96ec-d2fbd33acc71"] = &app.Event{
		ID:    "bf268c55-11e0-4efc-96ec-d2fbd33acc71",
		Title: "ev2",
		Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local),
		End:   time.Date(2023, 1, 1, 12, 30, 0, 0, time.Local),
	}
	storage.eventsByID["9ed65d93-987a-4f2f-9570-134648d192ae"] = &app.Event{
		ID:    "9ed65d93-987a-4f2f-9570-134648d192ae",
		Title: "ev3",
		Start: time.Date(2023, 1, 5, 12, 0, 0, 0, time.Local),
		End:   time.Date(2023, 1, 5, 12, 30, 0, 0, time.Local),
	}
	storage.eventsByID["02139778-5af7-4f90-a18f-8dd7a840d4cb"] = &app.Event{
		ID:    "02139778-5af7-4f90-a18f-8dd7a840d4cb",
		Title: "ev4",
		Start: time.Date(2023, 1, 29, 12, 0, 0, 0, time.Local),
		End:   time.Date(2023, 1, 29, 12, 30, 0, 0, time.Local),
	}

	ee, err := storage.GetEventsForDay(context.Background(), time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 0, len(ee))

	ee, err = storage.GetEventsForDay(context.Background(), time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 2, len(ee))

	ee, err = storage.GetEventsForWeek(context.Background(), time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 3, len(ee))

	ee, err = storage.GetEventsForMonth(context.Background(), time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 4, len(ee))
}

func TestMemStorage_UpdateEvent(t *testing.T) {
	storage, err := newMemStorage(context.Background())
	require.NoError(t, err, "unexpected error during mem storage init")

	ev1Id := "ebe62a7c-ce15-4f42-b34f-650c947f933f"
	ev1Start := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)
	ev1End := ev1Start.Add(15 * time.Minute)
	ev1UpdTitle := "upd ev1"

	ev2Id := "a6e61289-4ff2-43c9-aa0c-12037b0c4d41"
	ev2Start := time.Date(2023, 1, 2, 12, 0, 0, 0, time.Local)
	ev2End := ev2Start.Add(15 * time.Minute)

	storage.eventsByID[ev1Id] = &app.Event{ID: ev1Id, Title: "ev1", Start: ev1Start, End: ev1End}
	storage.eventsByID[ev2Id] = &app.Event{ID: ev2Id, Title: "ev2", Start: ev2Start, End: ev2End}

	testCases := []struct {
		name        string
		event       *app.Event
		expectedErr error
	}{
		{
			name:        "not found err",
			event:       &app.Event{ID: "not-found", Title: ev1UpdTitle, Start: ev1Start, End: ev1End},
			expectedErr: errNotFound,
		},
		{
			name:        "blank title err",
			event:       &app.Event{ID: ev1Id, Title: ""},
			expectedErr: errBlankTitle,
		},
		{
			name:        "blank start time",
			event:       &app.Event{ID: ev1Id, Title: ev1UpdTitle},
			expectedErr: errBlankStart,
		},
		{
			name:        "blank end time",
			event:       &app.Event{ID: ev1Id, Title: ev1UpdTitle, Start: ev1Start},
			expectedErr: errBlankEnd,
		},
		{
			name:        "reserved time err",
			event:       &app.Event{ID: ev1Id, Title: ev1UpdTitle, Start: ev2Start, End: ev2End},
			expectedErr: errReservedTime,
		},
		{
			name:  "successful update",
			event: &app.Event{ID: ev1Id, Title: ev1UpdTitle, Start: ev1Start, End: ev1End},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := storage.UpdateEvent(context.Background(), tc.event)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, 2, len(storage.eventsByID))
			} else {
				require.NoError(t, err, "unexpected error")
				require.Equal(t, 2, len(storage.eventsByID))
				require.Equal(t, ev1UpdTitle, storage.eventsByID[ev1Id].Title)
			}
		})
	}
}

func TestMemStorage_DeleteEvent(t *testing.T) {
	storage, err := newMemStorage(context.Background())
	require.NoError(t, err, "unexpected error during mem storage init")
	existedID := "0d84eec6-e3eb-4679-9ff7-0d57d6aa6bae"
	storage.eventsByID[existedID] = &app.Event{ID: existedID}

	testCases := []struct {
		name        string
		id          string
		expectedErr error
	}{
		{name: "failed delete", id: "n0tf0und-ffc5-463d-8b8f-6cdfbd70053d", expectedErr: errNotFound},
		{name: "successful delete", id: existedID, expectedErr: nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := storage.DeleteEvent(context.Background(), tc.id)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, 1, len(storage.eventsByID))
			} else {
				require.NoError(t, err, "unexpected error")
				require.Equal(t, 0, len(storage.eventsByID))
			}
		})
	}
}
