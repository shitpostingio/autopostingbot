package manager

import "time"

// GetPostingRate returns the current posting rate
func GetPostingRate() time.Duration {
	return manager.postingRate
}

// GetNextPostTime returns the time at which the next post is scheduled
func GetNextPostTime() time.Time {
	return manager.nextPostScheduled
}
