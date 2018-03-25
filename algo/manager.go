package algo

import (
	"fmt"
	"time"

	"gitlab.com/shitposting/autoposting-bot/database/entities"

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
	db               *gorm.DB
	AddChannel       chan tgbotapi.Update
	hourlyPostSignal <-chan time.Time
	hourlyPostRate   time.Duration
	postSignal       <-chan time.Time
}

// ManagerConfig is the configuration wanted for a given Manager instance.
// While BotAPIInstance is necessary, DatabasePath is not: if not present,
// Manager will try to load an existing database from ./autopostingbot.db,
// as per config.go.
type ManagerConfig struct {
	DatabasePath   string
	BotAPIInstance *tgbotapi.BotAPI
}

// NewManager returns a new Manager instance
func NewManager(mc ManagerConfig) (m Manager, err error) {
	m = Manager{
		botAPI:     mc.BotAPIInstance,
		AddChannel: make(chan tgbotapi.Update),
	}

	// Initialize gorm
	m.db, err = gorm.Open("sqlite3", mc.DatabasePath)
	if err != nil {
		return
	}

	// TODO: invoke algorithm to check when we'll have to post another photo,
	// and initialize postSignal with the output of time.After()

	// Calculate the hourly post rate on the current post availability
	m.calculateHourlyPostRate()

	// Initialize the calculation signal
	m.hourlyPostSignal = time.After(1 * time.Hour)

	// Initialize the postSignal on the hourlyRate
	if m.hourlyPostRate != 0 {
		m.postSignal = time.After(m.hourlyPostRate * time.Minute)
	}

	// Start the manager lifecycle
	go m.managerLifecycle()

	return
}

// managerLifecycle is the function we run indefinitely in a goroutine.
// It handles incoming updates, and the posting routine.
func (m *Manager) managerLifecycle() {
	for {
		select {
		case newUpdate := <-m.AddChannel:
			fmt.Println(newUpdate)
		case <-m.postSignal:
			fmt.Println("gotta post!")
			// reinitialize the posting signal with the
			m.postSignal = time.After(m.hourlyPostRate)
		case <-m.hourlyPostSignal:
			// calculate the new hourly post rate
			m.calculateHourlyPostRate()

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
		m.hourlyPostRate = time.Duration(60 / ppH)
		return
	}

	m.hourlyPostRate = 0
}
