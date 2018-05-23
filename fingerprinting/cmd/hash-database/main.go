package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/loglog/loglogclient"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/fingerprinting"

	"github.com/jinzhu/gorm"
)

var (
	configFilePath string

	// Importing loglog client
	l *loglogclient.LoglogClient
)

func main() {
	flag.StringVar(&configFilePath, "name", "./config.toml", "autoposting-bot configuration file")
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

	db, err := gorm.Open("mysql", cfg.DatabaseConnectionString())
	if err != nil {
		l.Err(err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		l.Err(err.Error())
	}

	defer db.Close()
	var duplicate []entities.Post
	db.Find(&duplicate)
	fmt.Println(duplicate)

	tx := db.Begin()
	for _, post := range duplicate {
		if !post.IsImage(db) || post.MediaHash != "" {
			continue
		}
		hash, err := fingerprinting.GetPhotoFingerprint(bot, post.Media)
		if err != nil {
			l.Err(fmt.Sprintf("cannot calculate fingerprint for media with ID %s", post.Media))
			continue
		}

		post.MediaHash = hash
		tx.Save(&post)
	}
	tx.Commit()
}
