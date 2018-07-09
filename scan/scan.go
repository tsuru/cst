package scan

import "time"

// Status is a type used to indicate the current state of an analysis.
type Status string

const (
	// StatusAborted indicates scan was aborted.
	StatusAborted = Status("aborted")

	// StatusFinished indicates scan was finished.
	StatusFinished = Status("finished")

	// StatusRunning indicates scan is running.
	StatusRunning = Status("running")

	// StatusScheduled indicates scan was scheduled (isn't running yet).
	StatusScheduled = Status("scheduled")
)

// Scan represents an analysis request over several security Scanners.
type Scan struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	Status     Status    `bson:"status,omitempty" json:"status"`
	Image      string    `bson:"image,omitempty" json:"image"`
	CreatedAt  time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	FinishedAt time.Time `bson:"finishedAt,omitempty" json:"finishedAt,omitempty"`
	Result     []Result  `bson:"result,omitempty" json:"result,omitempty"`
}

// Result holds an analysis result reported by a specific security scanner.
type Result struct {
	Scanner         string      `bson:"scanner" json:"scanner"`
	Vulnerabilities interface{} `bson:"vulnerabilities,omitempty" json:"vulnerabilities,omitempty"`
	Error           string      `bson:"error,omitempty" json:"error,omitempty"`
}
