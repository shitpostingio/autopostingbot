package manager

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	configuration "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

// Manager contains a copy of the needed fields and posting related values
type Manager struct {
	bot         *tgbotapi.BotAPI
	db          *gorm.DB
	config      *configuration.Config
	isDebugging bool
	isTesting   bool

	/* POSTING */
	nextPostScheduled   time.Time
	previousPostTime    time.Time
	previousPauseTime   time.Time
	postingRate         time.Duration
	postingSignal       <-chan time.Time
	postingRateChanging chan bool
	postingRateChanged  chan bool
}

const (
	lowPostsAlert     = "🚨 We're running out of memes!\nPosts enqueued: %d"
	imageTypeName     = "image"
	videoTypeName     = "video"
	animationTypeName = "animation"
)

var manager Manager

// StartManager starts the manager
func StartManager(botAPI *tgbotapi.BotAPI, database *gorm.DB, configuration *configuration.Config, debug, testing bool) error {

	manager = Manager{
		bot:         botAPI,
		db:          database,
		config:      configuration,
		isDebugging: debug,
		isTesting:   testing,
	}

	/* CACHE THE TYPES OF MEDIA */
	err := cacheTypes(database)
	if err != nil {
		return err
	}

	/* CREATE CHANNELS */
	manager.postingRateChanging = make(chan bool)
	manager.postingRateChanged = make(chan bool)

	/* SET THE HOURLY POSTING RATE */
	CalculateRateAndSchedulePosting(false)

	/* SET VALUES */
	manager.previousPauseTime = time.Now().Add(-10 * time.Minute)
	manager.previousPostTime = time.Now().Add(-10 * time.Minute)

	/* START THE MANAGER LIFECYCLE */
	go managePosting()

	return nil
}

// managePosting handles the various posting signals
func managePosting() {
	for {
		select {
		case <-manager.postingSignal:
			err := post()
			if err != nil {
				log.Error(err.Error())
			}
		case <-manager.postingRateChanging:
			<-manager.postingRateChanged
		}
	}
}

// cacheTypes caches the media types
func cacheTypes(db *gorm.DB) error {

	imageType := database.GetTypeByName(imageTypeName, db)
	types.Image = imageType.ID
	if imageType.ID == 0 {
		return errors.New("missing image type")
	}

	videoType := database.GetTypeByName(videoTypeName, db)
	types.Video = videoType.ID
	if videoType.ID == 0 {
		return errors.New("missing video type")
	}

	animationType := database.GetTypeByName(animationTypeName, db)
	types.Animation = animationType.ID
	if animationType.ID == 0 {
		return errors.New("missing animation type")
	}

	return nil
}

// notifyIfTooFewPosts sends a broadcast message to all the users in case posts queue is below a certain threshold
func notifyIfTooFewPosts(posts int) {

	if posts < manager.config.PostAlertThreshold {
		for _, user := range database.GetAllUsers(manager.db) {
			_, _ = manager.bot.Send(tgbotapi.NewMessage(int64(user.TelegramID), fmt.Sprintf(lowPostsAlert, posts)))
		}
	}

}
