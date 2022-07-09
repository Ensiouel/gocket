package gocket

import "encoding/json"

const (
	EventTypeEmit = "emit"
	EventTypePong = "pong"
)

type EventData []byte

func (d EventData) Unmarshal(v any) error {
	return json.Unmarshal(d, &v)
}

func (d EventData) String() string {
	return string(d)
}

type EventFunc func(data EventData)

type EventRequest struct {
	Type string      `json:"type"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type EventResponse struct {
	Type string          `json:"type"`
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

func (e EventResponse) IsEmit() bool {
	return e.Type == EventTypeEmit
}

func (e EventResponse) IsPong() bool {
	return e.Type == EventTypePong
}
