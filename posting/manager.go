package posting

import (
	log "github.com/sirupsen/logrus"
	configuration "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/posting/edition"
	"time"
)

type Manager struct {
	config      *configuration.Config
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
	requestPostChannel  chan *entities.Post
	requestPauseChannel chan time.Duration
	timedPostChannel    chan bool
}

var (
	m        Manager
	editions = map[string]edition.Edition{
		"shitposting": edition.ShitpostEdition{},
		"sushiporn":   edition.SushiPornEdition{},
	}
)

func Start(config *configuration.Config, debug, testing bool) {

	m.config = config
	m.isDebugging = debug
	m.isTesting = testing
	//TODO: LEVARE STRINGA CABLATA
	var found bool
	m.e, found = editions["shitposting"]
	if !found {
		log.Fatal("edition not found")
	}

	//

	m.requestPostChannel = make(chan *entities.Post)
	m.requestPauseChannel = make(chan time.Duration)
	m.timer = time.NewTimer(time.Minute)

	schedulePosting(time.Unix(0, 0))

}

func Listen() {

	var err error

	for {
		select {

		case postRequest := <-m.requestPostChannel:
			err = tryPosting(postRequest)
		case pauseRequest := <-m.requestPauseChannel:
			err = tryPausing(pauseRequest)
		case <-m.timer.C:
			err = postScheduled()
		}

		if err != nil {
			log.Error(err)
		}

	}

}
