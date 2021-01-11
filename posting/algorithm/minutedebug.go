package algorithm

import (
	"time"
)

// MinuteDebugAlgorithm is a debug structure to test posting functionalities.
type MinuteDebugAlgorithm struct{}

// GetNewPostingRate returns the new posting rate according to the algorithm
// and the queue length.
func (a MinuteDebugAlgorithm) GetNewPostingRate(queueLength int) time.Duration {
	return a.timeToNextPost(queueLength)
}

// timeToNextPost returns the estimate time until the next post, given
// the queue length.
func (MinuteDebugAlgorithm) timeToNextPost(queueLength int) time.Duration {

	if queueLength == 0 {
		return 0
	}

	return time.Minute

}

// EstimatePostTime estimates the amount of time that will pass before
// being able to post a certain media.
func (a MinuteDebugAlgorithm) EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	for i := queueLength; i > 0; i-- {
		totalDuration += a.timeToNextPost(i)
	}

	return

}

// GetAlgorithmName returns the algorithm name
func (MinuteDebugAlgorithm) GetAlgorithmName() string {
	return "debug"
}
