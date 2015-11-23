package twitter

import "time"

// Tweet struct
type Tweet struct {
	Success   bool
	ID        string
	UserID    string
	UserName  string
	Text      string
	CreatedAt time.Time
}

// FetchByID fetch tweet by self id
func (tweet *Tweet) FetchByID() (err error) {
	return
}
