package persistence

import "time"

type LogEntry struct {
	RequestID     string    `dynamodbav:"RequestID",json:"request_id"`
	UnixTimestamp int64     `dynamodbav:"UnixTimestamp",json:"unix_timestamp"`
	Timestamp     time.Time `dynamodbav:"Timestamp",json:"timestamp"`
	SearchQuery   string    `dynamodbav:"SearchQuery",json:"search_query"`
	ActualTitle   string    `dynamodbav:"ActualTitle",json:"actual_title"`
	Locale        string    `dynamodbav:"Locale",json:"locale"`
	UserID        string    `dynamodbav:"UserID",json:"user_id"`
	SessionID     string    `dynamodbav:"SessionID",json:"session_id"`
}

type Persistence interface {
	LogDefineIntentRequest(LogEntry)
}
