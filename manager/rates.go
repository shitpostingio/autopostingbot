package manager

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/database/database"
)

// CalculateRateAndSchedulePosting calculates the hourly post rate
// and saves it in the `Manager` instance.
func CalculateRateAndSchedulePosting(sendChangeSignals bool) {

	if sendChangeSignals {
		manager.postingRateChanging <- true
		defer func() { manager.postingRateChanged <- true }()
	}

	queueLength := database.GetQueueLength(manager.db)
	notifyIfTooFewPosts(queueLength)

	algo.UpdatePostingRate(&manager.postingRate, queueLength)
	if manager.postingRate == 0 {
		return
	}

	manager.nextPostScheduled = time.Now().Add(manager.postingRate)
	manager.postingSignal = time.After(manager.postingRate)
	return
}

// PausePosting pauses the posting for at least one hour
func PausePosting(howMuch time.Duration) error {

	/* PAUSES SHOULD NOT BE TOO FREQUENT */
	timeFromLastPause := -time.Until(manager.previousPauseTime)
	if timeFromLastPause < 5*time.Minute {

		log.Warn(fmt.Sprintf(
			"Pausing the posting rate too frequently is not allowed. Last pause was %s ago",
			timeFromLastPause))

		if !manager.isTesting {
			return errors.New("the previous pause was less than 5 minutes ago")
		}

		log.Info("Proceeding anyway because we're testing")
	}

	/* FORCE THE MANAGER TO STOP WAITING FOR THE NEXT POSTING SIGNAL */
	manager.postingRateChanging <- true

	/* DELAY BY AT LEAST ONE HOUR */
	if howMuch < time.Hour {
		howMuch = time.Hour
	}

	/* UPDATE MANAGER */
	if manager.isTesting {
		howMuch = 20 * time.Second
		manager.nextPostScheduled = time.Now()
		manager.postingRate = 0
	}

	manager.postingRate += howMuch
	manager.nextPostScheduled = manager.nextPostScheduled.Add(howMuch)
	manager.postingSignal = time.After(time.Until(manager.nextPostScheduled))
	manager.previousPauseTime = time.Now()

	/* ALLOW THE MANAGER TO GO ON */
	manager.postingRateChanged <- true
	return nil
}
