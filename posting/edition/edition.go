package edition

import (
	"time"
)

type Edition interface {
	GetNewPostingRate(queueLength int) time.Duration
	EstimatePostTime(queueLength int) time.Duration
	GetEditionName() string
}
