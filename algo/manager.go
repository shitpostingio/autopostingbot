package algo

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

// Manager is the central point of input/output for @AntonioBusillo's algorithm.
// It handles:
//  - channel updates
//  - database updates
//  - algorithm lifecycle
type Manager struct {
	botAPI           *tgbotapi.BotAPI
	channelID        int64
	db               *gorm.DB
	AddChannel       chan entities.Post
	hourlyPostSignal <-chan time.Time
	hourlyPostRate   time.Duration
	postSignal       <-chan time.Time
}

// ManagerConfig is the configuration wanted for a given Manager instance.
// While BotAPIInstance is necessary, DatabasePath is not: if not present,
// Manager will try to load an existing database from ./autopostingbot.db,
// as per config.go.
type ManagerConfig struct {
	DatabaseString string
	BotAPIInstance *tgbotapi.BotAPI
	ChannelID      int64
}

// NewManager returns a new Manager instance
func NewManager(mc ManagerConfig) (m Manager, err error) {
	if mc.ChannelID == 0 {
		err = errors.New("ChannelID is empty")
		return
	}

	m = Manager{
		botAPI:     mc.BotAPIInstance,
		channelID:  mc.ChannelID,
		AddChannel: make(chan entities.Post),
	}

	// Initialize gorm
	m.db, err = gorm.Open("mysql", mc.DatabaseString)
	if err != nil {
		return
	}

	// TODO: invoke algorithm to check when we'll have to post another photo,
	// and initialize postSignal with the output of time.After()

	// Calculate the hourly post rate on the current post availability
	m.calculateHourlyPostRate()

	// Print the hourly posting rate in minutes
	utility.YellowLog("Initial hourly posting rate set to " + (m.hourlyPostRate * time.Minute).String())

	// Initialize the calculation signal
	m.hourlyPostSignal = time.After(1 * time.Hour)

	// Initialize the postSignal on the hourlyRate
	m.setUpPostSignal()

	// Start the manager lifecycle
	go m.managerLifecycle()

	return
}

// managerLifecycle is the function we run indefinitely in a goroutine.
// It handles incoming updates, and the posting routine.
func (m *Manager) managerLifecycle() {
	for {
		select {
		case newPost := <-m.AddChannel:
			utility.GreenLog("got a new media to add!")
			// find the user id based on the telegram id
			var user entities.User
			m.db.Where("telegram_id = ?", newPost.UserID).First(&user)

			// if no user with said telegram_id was found, throw an error
			if user.ID == 0 {
				utility.PrettyError(fmt.Errorf("cannot find user_id for telegram_id %d", newPost.UserID))
				continue
			}

			// set the entity id to the database's user id
			newPost.UserID = user.ID

			// add to the database
			m.db.Create(&newPost)
		case <-m.postSignal:
			utility.GreenLog("it's time to post!")
			wtp, err := m.whatToPost()
			if err != nil {
				utility.PrettyError(err)
				continue
			}

			m.popAndPost(wtp)
			utility.GreenLog("all done!")
			m.setUpPostSignal()
		case <-m.hourlyPostSignal:
			utility.YellowLog("calculating the hourly posting rate...")
			// calculate the new hourly post rate
			m.calculateHourlyPostRate()

			// set up the posting signal, even if we already did that before
			m.setUpPostSignal()

			// see you in an hour!
			m.hourlyPostSignal = time.After(1 * time.Hour)
		}
	}
}

// calculateHourlyPostRate calculate the hourly post rate, and saves it in the Manager
// instance.
func (m *Manager) calculateHourlyPostRate() {
	var postsQueue []entities.Post
	m.db.Find(&postsQueue)

	ppH := postsPerHour(postsQueue)
	if ppH > 0 {
		m.hourlyPostRate = time.Duration(60/ppH) * time.Second
		return
	}

	m.hourlyPostRate = 0
}

// setUpPostSignal sets up the posting signal if there's something to post
func (m *Manager) setUpPostSignal() {
	if m.hourlyPostRate != 0 {
		m.postSignal = time.After(m.hourlyPostRate * time.Minute)
	}
}

// whatToPost returns the oldest media in the queue
func (m *Manager) whatToPost() (entities.Post, error) {
	var postsQueue []entities.Post
	m.db.Find(&postsQueue)
	sort.Sort(entities.Posts(postsQueue))

	if len(postsQueue) <= 0 {
		return entities.Post{}, errors.New("no element to post has been found")
	}

	return postsQueue[0], nil

}

// popAndPost removes entity from the database and post it to the Telegram channel
// identified by Manager.channelID
func (m *Manager) popAndPost(entity entities.Post) error {
	caption := "@shitpost"
	if entity.Caption != "" {
		entity.Caption = strings.TrimSpace(strings.Split(entity.Caption, "@shitpost")[0])
		caption = fmt.Sprintf("%s\n@shitpost", entity.Caption)
	}

	photoConfig := tgbotapi.NewPhotoShare(m.channelID, entity.Media)
	photoConfig.Caption = caption

	_, err := m.botAPI.Send(photoConfig)

	// checking if there's an error here gives us the chance to remove the posted
	// entity if everything was ok - if it wasn't, the error will be handled on the caller
	// level.
	if err == nil {
		m.db.Delete(&entity)
	}
	return err
}
