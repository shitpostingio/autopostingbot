package main

import (
	"flag"

	"gitlab.com/shitposting/autoposting-bot/database/migrations"

	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// autoposting-deploy-db will create a database file, with all the tables ready to be used

var (
	dbName string
)

func main() {
	flag.StringVar(&dbName, "name", "autopostingbot.db", "database filename, will accept a complete absolute/relative path too")
	flag.Parse()

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		utility.PrettyFatal(err)
	}
	defer db.Close()

	// Create User table
	migrations.CreateUsers(db)
	// Create Category table
	migrations.CreateCategories(db)
	// Create Posts table
	migrations.CreatePosts(db)
}
