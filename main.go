package main

import (
	"flag"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
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

// TODO: this function will be removed as soon as we have some code that actually uses gorm
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func dumbORM() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}
