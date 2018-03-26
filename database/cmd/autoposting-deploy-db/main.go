package main

import (
	"flag"

	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/migrations"

	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// autoposting-deploy-db will create a database file, with all the tables ready to be used

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

	defer db.Close()
	// Create User table
	migrations.CreateUsers(db)
	// Create Category table
	migrations.CreateCategories(db)
	// Create Posts table
	migrations.CreatePosts(db)
}
