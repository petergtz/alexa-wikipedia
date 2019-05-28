package persistence

import "time"

type LogEntry struct {
	Timestamp   time.Time
	SearchQuery string
	ActualTitle string
	Locale      string
}

type Persistence interface {
	LogDefineIntentRequest(LogEntry)
}
