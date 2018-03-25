package main

import (
	"flag"
	"log"

	"gitlab.com/shitposting/autoposting-bot/database/entities"

	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// autoposting-add-user will add an approved user to the autoposting-bot database
var (
	dbName string
	userID int
)

func main() {
	flag.StringVar(&dbName, "name", "autopostingbot.db", "database filename, will accept a complete absolute/relative path too")
	flag.IntVar(&userID, "userid", 0, "ID of the user you want to add")
	flag.Parse()

	if userID == 0 {
		log.Fatal("user ID is needed")
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		utility.PrettyFatal(err)
	}
	defer db.Close()

	db.Create(&entities.User{TelegramID: userID})
}
