package main

import (
	"flag"
	"fmt"
	"net/http"

	"gitlab.com/shitposting/autoposting-bot/algo"
	cfg "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/entities"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/shitposting/autoposting-bot/command"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

var (
	// parsed config file
	config cfg.Config

	// config file path, if not specified it will read
	// ./config.toml
	configFilePath string

	// Version is the autoposting-bot version, a compile-time value
	Version string

	// Build is the git tag for the current version
	Build string

	manager algo.Manager

	db *gorm.DB

	err   error
	debug bool
)

func main() {
	setupCLIParams()

	utility.GreenLog(fmt.Sprintf("Shitposting autoposting-bot version %s, build %s\n", Version, Build))
	utility.YellowLog(fmt.Sprintf("INFO - reading configuration file located at %s", configFilePath))
	config, err = cfg.ReadConfigFile(configFilePath)
	if err != nil {
		utility.PrettyFatal(err)
	}

	// setup a Telegram bot API instance
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		utility.PrettyFatal(err)
	}

	// should we activate debug output?
	bot.Debug = debug

	//utility.YellowLog
	utility.YellowLog(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	// set webhook to an adequate value
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.WebHookURL()))
	if err != nil {
		utility.PrettyFatal(err)
	}

	updates := bot.ListenForWebhook(config.WebHookPath())

	manager, err = algo.NewManager(algo.ManagerConfig{
		ChannelID:      int64(config.ChannelID),
		BotAPIInstance: bot,
		DatabaseString: config.DatabaseConnectionString(),
		Debug:          debug,
	})

	if err != nil {
		utility.PrettyFatal(err)
	}

	// Initialize gorm
	db, err = gorm.Open("mysql", config.DatabaseConnectionString())
	if err != nil {
		utility.PrettyFatal(err)
	}

	go startServer()

	for update := range updates {
		go func(update tgbotapi.Update, bot *tgbotapi.BotAPI, manager algo.Manager) {
			if iCanUseThis(update) {
				err := command.Handle(update, bot, &manager)
				if err != nil {
					utility.PrettyError(err)
				}
			}
		}(update, bot, manager)
	}

}

func setupCLIParams() {
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.BoolVar(&debug, "debug", false, "activate all the debug features")
	flag.Parse()
}

func startServer() {
	if config.TLS {
		go utility.PrettyFatal(http.ListenAndServeTLS(config.BindString(), config.TLSCertPath, config.TLSKeyPath, nil))
	}

	go utility.PrettyFatal(http.ListenAndServe(config.BindString(), nil))
}

func iCanUseThis(update tgbotapi.Update) bool {
	realUpdate := &tgbotapi.Message{}
	if update.Message != nil && update.Message.From != nil {
		realUpdate = update.Message
	} else if update.EditedMessage != nil && update.EditedMessage.From != nil {
		realUpdate = update.EditedMessage
	} else {
		return false
	}

	destID := realUpdate.From.ID
	var users []entities.User

	db.Where("telegram_id = ?", destID).Find(&users)
	if len(users) > 0 {
		return true
	}

	return false
}
