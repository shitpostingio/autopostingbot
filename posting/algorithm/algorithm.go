package algorithm

import (
	"time"
)

// Algorithm is the interface used to implement posting strategies.
type Algorithm interface {

	// GetNewPostingRate returns the posting rate given the queue length.
	GetNewPostingRate(queueLength int) time.Duration

	// EstimatePostTime estimates the amount of time that will pass before
	// being able to post a certain media.
	EstimatePostTime(queueLength int) time.Duration

	// GetAlgorithmName returns the algorithm name
	GetAlgorithmName() string
}
