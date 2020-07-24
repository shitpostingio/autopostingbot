package edition

import (
	"math/rand"
	"time"
)

type SushiPornEdition struct {
}

func (SushiPornEdition) GetEditionName() string {
	return "sushiporn"
}

// GetNewPostingRate updates the posting rate to a new value based
// on the edition of the active edition
func (e SushiPornEdition) GetNewPostingRate(queueLength int) time.Duration {
	return e.timeToNextPost(queueLength)
}

func (SushiPornEdition) timeToNextPost(queueLength int) time.Duration {

	if queueLength == 0 {
		return 0
	}

	rand.Seed(time.Now().UnixNano())
	timeToWait := rand.Intn(4) + 4
	return time.Duration(timeToWait) * time.Hour

}

// postsPerHour will distribute the amount of posts in the queue
// over a 24 hour period
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

// EstimatePostTime estimates the time until the posting of a
// certain item in the queue
func (e SushiPornEdition) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += e.timeToNextPost(i)
	}

	return

}
