package posting

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/posting/edition"
	"time"
)

type RequestPostStruct struct {
	Post      *entities.Post
	ErrorChan chan error
}

type RequestPauseStruct struct {
	Duration  time.Duration
	ErrorChan chan error
}

func RequestPost(post *entities.Post) error {

	rps := RequestPostStruct{
		Post:      post,
		ErrorChan: make(chan error, 1),
	}

	m.requestPostChannel <- rps
	return <-rps.ErrorChan

}

func RequestPause(duration time.Duration) error {

	rps := RequestPauseStruct{
		Duration:  duration,
		ErrorChan: make(chan error, 1),
	}

	m.requestPauseChannel <- rps
	return <-rps.ErrorChan

}

// GetPostingRate returns the current posting rate
func GetPostingRate() time.Duration {
	return m.postingRate
}

// GetNextPostTime returns the time at which the next post is scheduled
func GetNextPostTime() time.Time {
	return m.nextPostScheduled
}

func GetPostingManager() edition.Edition {
	return m.e
}

func ForcePostScheduling() {
	schedulePosting(time.Unix(0, 0))
}
