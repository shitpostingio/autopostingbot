package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/shitposting/autoposting-bot/algo"
	cfg "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/loglog/loglogclient"

	"github.com/empetrone/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/shitposting/autoposting-bot/command"
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

	manager *algo.Manager

	db *gorm.DB

	err   error
	debug bool

	// Importing loglog client
	l *loglogclient.LoglogClient
)

func main() {
	setupCLIParams()

	config, err = cfg.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	l = loglogclient.NewClient(
		loglogclient.Config{
			SocketPath:    config.SocketPath,
			ApplicationID: "Autoposting-bot",
		})

	l.Info(fmt.Sprintf("Shitposting autoposting-bot version %s, build %s\n", Version, Build))
	l.Info(fmt.Sprintf("INFO - reading configuration file located at %s", configFilePath))

	// setup a Telegram bot API instance
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		l.Err(err.Error())
		os.Exit(1)
	}

	// should we activate debug output?
	bot.Debug = debug

	l.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	// set webhook to an adequate value
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.WebHookURL()))
	if err != nil {
		l.Err(err.Error())
		os.Exit(1)
	}

	updates := bot.ListenForWebhook(config.WebHookPath())

	manager, err = algo.NewManager(algo.ManagerConfig{
		ChannelID:      int64(config.ChannelID),
		BotAPIInstance: bot,
		DatabaseString: config.DatabaseConnectionString(),
		Debug:          debug,
		Log:            l,
	})

	if err != nil {
		l.Err(err.Error())
		os.Exit(1)
	}

	// Initialize gorm
	db, err = gorm.Open("mysql", config.DatabaseConnectionString())
	if err != nil {
		l.Err(err.Error())
		os.Exit(1)
	}

	go startServer()

	for update := range updates {
		go func(update tgbotapi.Update, bot *tgbotapi.BotAPI, manager *algo.Manager) {
			if iCanUseThis(update) {
				err := command.Handle(update, bot, manager)
				if err != nil {
					l.Err(err.Error())
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
		go l.Err((http.ListenAndServeTLS(config.BindString(), config.TLSCertPath, config.TLSKeyPath, nil)).Error())
	}

	go l.Err((http.ListenAndServe(config.BindString(), nil)).Error())
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
