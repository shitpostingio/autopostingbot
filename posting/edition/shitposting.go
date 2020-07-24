package edition

import (
	"time"
)

type ShitpostingEdition struct {

}

func (ShitpostingEdition) GetEditionName() string {
	return "shitposting"
}

// GetNewPostingRate updates the posting rate to a new value based
// on the edition of the active edition
func (e ShitpostingEdition) GetNewPostingRate(queueLength int) time.Duration {
	return e.timeToNextPost(queueLength)
}

func (e ShitpostingEdition) timeToNextPost(queueLength int) time.Duration {

	//
	postPerHour := e.postsPerHour(queueLength)

	//
	if postPerHour == 0 {
		return 0
	}

	//
	return time.Duration(60/postPerHour) * time.Minute

}

// postsPerHour will distribute the amount of posts in the queue
// over a 24 hour period
func (ShitpostingEdition) postsPerHour(queueLength int) int {

	if queueLength == 0 {
		return 0
	}

	postsPerHour := queueLength / 24

	if postsPerHour < 3 {
		return postsPerHour + 1
	}

	return postsPerHour

}

// EstimatePostTime estimates the time until the posting of a
// certain item in the queue
func (e ShitpostingEdition) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += e.timeToNextPost(i)
	}

	return

}

