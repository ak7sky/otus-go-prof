package app

import "time"

type Event struct {
	ID           string
	Title        string
	Start        time.Time     `db:"event_start"`
	End          time.Time     `db:"event_end"`
	Desc         string        `db:"descr"`
	NotifyBefore time.Duration `db:"notify_before"`
}
