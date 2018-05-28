package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/fingerprinting"
	"gitlab.com/shitposting/loglog/loglogclient"

	"github.com/empetrone/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
			ApplicationID: "hash-database",
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
		if !post.IsImage(db) {
			continue
		}

		log.Println("processing photo with ID", post.ID)
		aHash, pHash, err := fingerprinting.GetPhotoFingerprint(bot, post.Media)
		if err != nil {
			l.Err(fmt.Sprintf("cannot calculate fingerprint for media with ID %s", post.Media))
			continue
		}

		post.AHash = aHash
		post.PHash = pHash
		tx.Save(&post)
	}
	tx.Commit()
}
