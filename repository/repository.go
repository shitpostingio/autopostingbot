package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	Config *structs.Config
	Tdlib  *client.Client
	Me     *client.User
)

// Repository represents a container for common-use variables
type Repository struct {
	Bot    *tgbotapi.BotAPI
	Db     *gorm.DB
	Config *structs.Config
}

// SetVariables sets variables in the repository
func SetVariables(botAPI *tgbotapi.BotAPI, database *gorm.DB, config *structs.Config) *Repository {

	/* CREATE REPOSITORY */
	var repo Repository

	/* SET VARIABLES */
	repo.SetBot(botAPI)
	repo.SetDatabase(database)
	repo.SetConfig(config)

	Config = config

	return &repo
}

//TODO: LEVARE STA ROBA
// SetBot sets the bot in the repository
func (repo *Repository) SetBot(api *tgbotapi.BotAPI) {
	repo.Bot = api
}

// SetDatabase sets the database in the repository
func (repo *Repository) SetDatabase(db *gorm.DB) {
	repo.Db = db
}

// SetConfig sets the configuration in the repository
func (repo *Repository) SetConfig(cfg *structs.Config) {
	repo.Config = cfg
}
