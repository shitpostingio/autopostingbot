package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"github.com/zelenin/go-tdlib/client"

	configuration "gitlab.com/shitposting/autoposting-bot/config"
)

var (
	Bot    *tgbotapi.BotAPI
	Db     *gorm.DB
	Config *configuration.Config
	Tdlib  *client.Client
	Me     *client.User
)

// Repository represents a container for common-use variables
type Repository struct {
	Bot    *tgbotapi.BotAPI
	Db     *gorm.DB
	Config *configuration.Config
}

// SetVariables sets variables in the repository
func SetVariables(botAPI *tgbotapi.BotAPI, database *gorm.DB, config *configuration.Config) *Repository {

	/* CREATE REPOSITORY */
	var repo Repository

	/* SET VARIABLES */
	repo.SetBot(botAPI)
	repo.SetDatabase(database)
	repo.SetConfig(config)

	Bot = botAPI
	Db = database
	Config = config

	return &repo
}

// SetBot sets the bot in the repository
func (repo *Repository) SetBot(api *tgbotapi.BotAPI) {
	repo.Bot = api
}

// SetDatabase sets the database in the repository
func (repo *Repository) SetDatabase(db *gorm.DB) {
	repo.Db = db
}

// SetConfig sets the configuration in the repository
func (repo *Repository) SetConfig(cfg *configuration.Config) {
	repo.Config = cfg
}
