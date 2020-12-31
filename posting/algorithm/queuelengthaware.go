package algorithm

import (
	"time"
)

// QueueLengthAwareAlgorithm represents the Shitpost posting algorithm.
type QueueLengthAwareAlgorithm struct{}

// GetNewPostingRate returns the new posting rate according to the algorithm
// and the queue length.
func (a QueueLengthAwareAlgorithm) GetNewPostingRate(queueLength int) time.Duration {
	return a.timeToNextPost(queueLength)
}

// timeToNextPost returns the estimate time until the next post, given
// the queue length.
func (a QueueLengthAwareAlgorithm) timeToNextPost(queueLength int) time.Duration {

	//
	postPerHour := a.postsPerHour(queueLength)

	//
	if postPerHour == 0 {
		return 0
	}

	//
	return time.Duration(60/postPerHour) * time.Minute

}

// postsPerHour returns how many posts per hour there will be
// in the next 24 hours period, based on the queue length.
func (QueueLengthAwareAlgorithm) postsPerHour(queueLength int) int {

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
func (a QueueLengthAwareAlgorithm) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += a.timeToNextPost(i)
	}

	return

}

// GetAlgorithmName returns the algorithm name
func (QueueLengthAwareAlgorithm) GetAlgorithmName() string {
	return "queue length aware"
}
