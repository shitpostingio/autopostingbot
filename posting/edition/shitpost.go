package edition

import (
	"time"
)

// ShitpostEdition represents the Shitpost posting algorithm.
type ShitpostEdition struct {}

// GetEditionName returns the name of the edition.
func (ShitpostEdition) GetEditionName() string {
	return "shitpost"
}

// GetNewPostingRate returns the new posting rate according to the algorithm
// and the queue length.
func (e ShitpostEdition) GetNewPostingRate(queueLength int) time.Duration {
	return e.timeToNextPost(queueLength)
}

// timeToNextPost returns the estimate time until the next post, given
// the queue length.
func (e ShitpostEdition) timeToNextPost(queueLength int) time.Duration {

	//
	postPerHour := e.postsPerHour(queueLength)

	//
	if postPerHour == 0 {
		return 0
	}

	//
	return time.Duration(60/postPerHour) * time.Minute

}

// postsPerHour returns how many posts per hour there will be
// in the next 24 hours period, based on the queue length.
func (ShitpostEdition) postsPerHour(queueLength int) int {

	if queueLength == 0 {
		return 0
	}

	postsPerHour := queueLength / 24

	if postsPerHour < 3 {
		return postsPerHour + 1
	}

	return postsPerHour

}

// EstimatePostTime estimates the amount of time that will pass before
// being able to post a certain media.
func (e ShitpostEdition) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += e.timeToNextPost(i)
	}

	return

}
