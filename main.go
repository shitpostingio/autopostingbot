package main

import (
	"flag"
	"net/http"

	"gitlab.com/shitposting/autoposting-bot/algo"

	"github.com/fatih/color"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/shitposting/autoposting-bot/command"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

var (
	// parsed config file
	config Config

	// config file path, if not specified it will read
	// ./config.toml
	configFilePath string

	// Version is the autoposting-bot version, a compile-time value
	Version string

	// Build is the git tag for the current version
	Build string

	manager algo.Manager

	err   error
	debug bool
)

func main() {
	setupCLIParams()
	color.Green("Shitposting autoposting-bot version %s, build %s", Version, Build)
	color.Yellow("INFO - reading configuration file located at %s", configFilePath)
	config, err = ReadConfigFile(configFilePath)
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

	color.Yellow("Authorized on account %s", bot.Self.UserName)

	// set webhook to an adequate value
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.WebHookURL()))
	if err != nil {
		utility.PrettyFatal(err)
	}

	updates := bot.ListenForWebhook(config.WebHookPath())

	manager, err = algo.NewManager(algo.ManagerConfig{
		DatabasePath:   config.DatabasePath,
		ChannelID:      config.ChannelID,
		BotAPIInstance: bot,
	})

	if err != nil {
		utility.PrettyFatal(err)
	}

	go startServer()

	for update := range updates {
		err := command.Handle(update, bot)
		if err != nil {
			utility.PrettyError(err)
		}
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
