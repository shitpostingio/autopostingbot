package posting

import (
	"github.com/shitpostingio/autopostingbot/config/structs"
	"github.com/shitpostingio/autopostingbot/posting/edition"
	log "github.com/sirupsen/logrus"
	"time"
)

// Manager is the posting manager.
// It contains information on posting rates, scheduling and the running edition.
type Manager struct {
	config      *structs.Config
	isDebugging bool

	/* POSTING */
	nextPostScheduled time.Time
	previousPostTime  time.Time
	previousPauseTime time.Time
	postingRate       time.Duration

	//
	e edition.Edition

	//
	timer *time.Timer

	//
	requestPostChannel  chan RequestPostStruct
	requestPauseChannel chan RequestPauseStruct
}

var (
	m        Manager
	editions = map[string]edition.Edition{
		"shitpost":  edition.ShitpostEdition{},
		"sushiporn": edition.SushiPornEdition{},
	}
)

// Start sets the Manager up and starts the post scheduling.
func Start(config *structs.Config, debug bool) {

	//
	m.config = config
	m.isDebugging = debug

	//
	var found bool
	m.e, found = editions[config.Autoposting.Edition]
	if !found {
		log.Fatal("edition not found")
	}

	//
	m.requestPostChannel = make(chan RequestPostStruct)
	m.requestPauseChannel = make(chan RequestPauseStruct)
	m.timer = time.NewTimer(time.Minute)

	//
	ForcePostScheduling()

}

// Listen makes the Manager listen for posting and pause requests.
func Listen() {

	var err error

	for {
		select {

		case postRequest := <-m.requestPostChannel:
			err = tryPosting(postRequest.Post)
			postRequest.ErrorChan <- err
		case pauseRequest := <-m.requestPauseChannel:
			err = tryPausing(pauseRequest.Duration)
			pauseRequest.ErrorChan <- err
		case <-m.timer.C:
			err = postScheduled()
		}

		if err != nil {
			log.Error(err)
		}

	}

}
