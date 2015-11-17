package trm

import (
	"fmt"
	"time"
)

const (
	OpRequest  = "Request"
	OpNetwork  = "Network"
	OpLimit    = "Limit"
	OpResponse = "Response"
)

// OpRequest implys Request Operation

// A TwitterError records a failed get of tweet.
type TwitterError struct {
	Op    string    // the failing Operation (Request, Network, Limit, Response)
	ID    string    // the twitter id
	Reset time.Time // the reset time
	Err   error     // the reason the get failed
}

func (e *TwitterError) Error() string {
	if e.Op != "Limit" {
		return e.Op + " https://twitter.com/statuses/" + e.ID + ": " + e.Err.Error()
	}
	s := fmt.Sprintf("reset at %v", e.Reset)
	return e.Op + " https://twitter.com/statuses/" + e.ID + ": " + s
}
