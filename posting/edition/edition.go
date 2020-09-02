package edition

import (
	"time"
)

// Edition is the interface used to implement posting strategies.
type Edition interface {

	// GetNewPostingRate returns the posting rate given the queue length.
	GetNewPostingRate(queueLength int) time.Duration

	// EstimatePostTime estimates the amount of time that will pass before
	// being able to post a certain media.
	EstimatePostTime(queueLength int) time.Duration

	// GetEditionName returns the name of the edition.
	GetEditionName() string
}
