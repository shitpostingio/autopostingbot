package main

import (
	"flag"
	"fmt"

	"gitlab.com/shitposting/autoposting-bot/database/entities"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/fingerprinting"

	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "name", "./config.toml", "autoposting-bot configuration file")
	flag.Parse()

	cfg, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		utility.PrettyFatal(err)
	}

	db, err := gorm.Open("mysql", cfg.DatabaseConnectionString())
	if err != nil {
		utility.PrettyFatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		utility.PrettyFatal(err)
	}

	defer db.Close()
	var duplicate []entities.Post
	db.Where("image_hash = ?", "").Find(&duplicate)

	tx := db.Begin()
	for _, post := range duplicate {
		if !post.IsImage(db) {
			continue
		}
		hash, err := fingerprinting.GetPhotoFingerprint(bot, post.Media)
		if err != nil {
			utility.PrettyError(fmt.Errorf("cannot calculate fingerprint for media with ID %s", post.Media))
			continue
		}

		post.ImageHash = hash
		tx.Save(&post)
	}
	tx.Commit()
}
