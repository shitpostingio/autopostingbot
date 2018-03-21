package algo

import (
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

// Manager is the central point of input/output for @AntonioBusillo's algorithm.
// It handles:
//  - channel updates
//  - database updates
//  - algorithm lifecycle
type Manager struct {
	botAPI     *tgbotapi.BotAPI
	db         *gorm.DB
	AddChannel chan tgbotapi.Update
	postSignal chan time.Time
}

// ManagerConfig is the configuration wanted for a given Manager instance.
// While BotAPIInstance is necessary, DatabasePath is not: if not present,
// Manager will try to load an existing database from ./autopostingbot.db.
type ManagerConfig struct {
	DatabasePath   string
	BotAPIInstance *tgbotapi.BotAPI
}

// NewManager returns a new Manager instance
func NewManager(mc ManagerConfig) (m Manager, err error) {
	if mc.DatabasePath == "" {
		mc.DatabasePath = "./autopostingbot.db"
	}

	m = Manager{
		botAPI:     mc.BotAPIInstance,
		AddChannel: make(chan tgbotapi.Update),
	}

	// TODO: invoke algorithm to check when we'll have to post another photo,
	// and initialize postSignal with the output of time.After()

	// Initialize gorm
	m.db, err = gorm.Open("sqlite3", mc.DatabasePath)

	// Start the manager lifecycle
	go m.managerLifecycle()

	return
}

// managerLifecycle is the function we run indefinitely in a goroutine.
// It handles incoming updates, and the posting routine.
func (m Manager) managerLifecycle() {
	for {
		select {
		case newUpdate := <-m.AddChannel:
			fmt.Println(newUpdate)
		case <-m.postSignal:
			fmt.Println("gotta post!")
		}
	}
}
