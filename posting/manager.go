package posting

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
	"gitlab.com/shitposting/autoposting-bot/posting/edition"
	"time"
)

type Manager struct {
	config      *structs.Config
	isDebugging bool
	isTesting   bool

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

func Start(config *structs.Config, debug, testing bool) {

	//
	m.config = config
	m.isDebugging = debug
	m.isTesting = testing

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

	schedulePosting(time.Unix(0, 0))

}

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
