package posting

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/posting/edition"
	"time"
)

func RequestPost(post *entities.Post) {
	m.requestPostChannel <- post
}

func RequestPause(duration time.Duration) {
	m.requestPauseChannel <- duration
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
