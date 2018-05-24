package loglogclient

import "time"

// Payload is what the loglog server handles, and transforms into a LogEntry database entity
type Payload struct {
	ApplicationID string
	Type          LogType
	Payload       string
	CreatedAt     time.Time
}
