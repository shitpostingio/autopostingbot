package main

import (
	"flag"
	"fmt"
	"gitlab.com/shitposting/autoposting-bot/config/old"
	"log"

	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/shitposting/loglog-ng"
)

var (
	userID         int
	userHandle     string
	configFilePath string
)

func main() {

	/* PARSE CLI ARGUMENTS */
	flag.StringVar(&configFilePath, "config", "./config.toml", "autoposting-bot configuration file")
	flag.IntVar(&userID, "userid", 0, "ID of the user you want to add")
	flag.StringVar(&userHandle, "handle", "", "Handle of the user you want to add")
	flag.Parse()

	/* END IF NO USER ID SPECIFIED */
	if userID == 0 {
		log.Fatal("must specify userID with the -userid flag")
	}

	/* LOAD CONFIG */
	cfg, err := old.Load(configFilePath, false)
	if err != nil {
		log.Fatal(err)
	}

	/* CONNECT TO DB */
	db, err := gorm.Open("mysql", cfg.DB.DatabaseConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(db.Close())
		}
	}()

	/* INSTANTIATE LOGLOG */
	err = loglog.Setup(cfg.LogLog.ApplicationID)
	if err != nil {
		log.Println("Loglog error:", err)
	}

	/* INSERT USER IN THE DATABASE */
	result := db.Create(&entities.User{TelegramID: userID, Handle: userHandle})
	if result.RowsAffected == 1 {
		loglog.Info(fmt.Sprintf("User with ID %d is now allowed to interact with the bot", userID))
	} else {
		log.Fatal(result.Error)
	}
}
