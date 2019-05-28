package persistence

import "time"

type LogEntry struct {
	UnixTimestamp int64     `json:"unix_timestamp"`
	Timestamp     time.Time `json:"timestamp"`
	SearchQuery   string    `json:"search_query"`
	ActualTitle   string    `json:"actual_title"`
	Locale        string    `json:"locale"`
	UserID        string    `json:"user_id"`
	SessionID     string    `json:"session_id"`
}

type Persistence interface {
	LogDefineIntentRequest(LogEntry)
}
