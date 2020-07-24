package algo

import (
	"math/rand"
	"time"

	"gitlab.com/shitposting/autoposting-bot/edition"
)

// GetNewPostingRate updates the posting rate to a new value based
// on the edition of the active edition
func UpdatePostingRate(postingRate *time.Duration, queueLength int) {

	switch {
	case edition.IsShitpost():
		useShitpostScheduling(postingRate, queueLength)
	case edition.IsSushiporn():
		useSushipornScheduling(postingRate, queueLength)
	}

}

// useSushipornScheduling uses the sushiporn edition
// to update the posting rate
func useSushipornScheduling(postingRate *time.Duration, queueLength int) {

	if queueLength == 0 {
		*postingRate = 0
		return
	}

	rand.Seed(time.Now().UnixNano())
	timeToWait := rand.Intn(4) + 4
	*postingRate = time.Duration(timeToWait) * time.Hour

}

// useShitpostScheduling uses the shitpost edition
// to update the posting rate
func useShitpostScheduling(postingRate *time.Duration, queueLength int) {

	postPerHour := postsPerHour(queueLength)

	if postPerHour == 0 {
		*postingRate = 0
		return
	}

	*postingRate = time.Duration(60/postPerHour) * time.Minute

}

// postsPerHour will distribute the amount of posts in the queue
// over a 24 hour period
func postsPerHour(queueLength int) int {

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
func EstimatePostTime(queueLength int) (totalDuration time.Duration) {

	var scheduleFunction func(*time.Duration, int)

	switch {
	case edition.IsShitpost():
		scheduleFunction = useShitpostScheduling
	case edition.IsSushiporn():
		scheduleFunction = useSushipornScheduling
	}

	if scheduleFunction == nil {
		return 0
	}

	var untilNextPost time.Duration
	for i := queueLength; i > 0; i-- {
		scheduleFunction(&untilNextPost, i)
		totalDuration += untilNextPost
	}

	return
}
