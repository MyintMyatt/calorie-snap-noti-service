package event

import "time"

type Event struct {
	ID string `json:"event_id"`
	Type string `json:"event_type"`
	Version int `json:"version"`
	Source string `json:"source"`
	ClientTime time.Time `json:"client_time"` // client create date time
	Data any `json:"data"`
}