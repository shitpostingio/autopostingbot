package main

import (
	"flag"
	"log"

	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/loglog/loglogclient"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// autoposting-add-user will add an approved user to the autoposting-bot database
var (
	userID         int
	configFilePath string

	// Importing loglog client
	l *loglogclient.LoglogClient
)

func main() {
	flag.StringVar(&configFilePath, "name", "./config.toml", "autoposting-bot configuration file")
	flag.IntVar(&userID, "userid", 0, "ID of the user you want to add")
	flag.Parse()

	cfg, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	l = loglogclient.NewClient(
		loglogclient.Config{
			SocketPath:    cfg.SocketPath,
			ApplicationID: "Autoposting-bot",
		})

	if userID == 0 {
		l.Warn("user ID is needed")
	}

	db, err := gorm.Open("mysql", cfg.DatabaseConnectionString())
	if err != nil {
		l.Err(err.Error())
	}
	defer db.Close()

	db.Create(&entities.User{TelegramID: userID})
}
