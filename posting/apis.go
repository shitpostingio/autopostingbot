package posting

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/posting/edition"
	"time"
)

// RequestPostStruct represents a request to post a entities.Post.
// It includes a channel to retrieve the posting result.
type RequestPostStruct struct {

	// Post is the media to be posted.
	Post      *entities.Post

	// ErrorChan is a channel to retrieve post results.
	ErrorChan chan error
}

// RequestPauseStruct represents a request to pause posting.
// It includes a channel to retrieve the operation result.
type RequestPauseStruct struct {

	// Duration is the amount of time to pause posting for.
	Duration  time.Duration

	// ErrorChan is a channel to the operation results.
	ErrorChan chan error
}

// RequestPost requests the posting of a media.
func RequestPost(post *entities.Post) error {

	rps := RequestPostStruct{
		Post:      post,
		ErrorChan: make(chan error, 1),
	}

	m.requestPostChannel <- rps
	return <-rps.ErrorChan

}

// RequestPause requests a pause in the posting.
func RequestPause(duration time.Duration) error {

	rps := RequestPauseStruct{
		Duration:  duration,
		ErrorChan: make(chan error, 1),
	}

	m.requestPauseChannel <- rps
	return <-rps.ErrorChan

}

// GetPostingRate returns the current posting rate.
func GetPostingRate() time.Duration {
	return m.postingRate
}

// GetNextPostTime returns the time at which the next post is scheduled.
func GetNextPostTime() time.Time {
	return m.nextPostScheduled
}

// GetPostingManager returns the current posting manager edition.
func GetPostingManager() edition.Edition {
	return m.e
}

// ForcePostScheduling forces a new post scheduling.
func ForcePostScheduling() {
	schedulePosting(time.Unix(0, 0))
}
