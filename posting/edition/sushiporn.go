package edition

import (
	"math/rand"
	"time"
)

// SushiPornEdition represents the Shitpost posting algorithm.
type SushiPornEdition struct{}

// GetEditionName returns the name of the edition.
func (SushiPornEdition) GetEditionName() string {
	return "sushiporn"
}

// GetNewPostingRate returns the new posting rate according to the algorithm
// and the queue length.
func (e SushiPornEdition) GetNewPostingRate(queueLength int) time.Duration {
	return e.timeToNextPost(queueLength)
}

// timeToNextPost returns the estimate time until the next post, given
// the queue length.
func (SushiPornEdition) timeToNextPost(queueLength int) time.Duration {

	if queueLength == 0 {
		return 0
	}

	rand.Seed(time.Now().UnixNano())
	timeToWait := rand.Intn(4) + 4
	return time.Duration(timeToWait) * time.Hour

}

// postsPerHour returns how many posts per hour there will be
// in the next 24 hours period, based on the queue length.
func (SushiPornEdition) postsPerHour(queueLength int) int {

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
func (e SushiPornEdition) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += e.timeToNextPost(i)
	}

	return

}
