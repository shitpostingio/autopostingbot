package main

import (
	"flag"
	"log"

	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/loglog/loglogclient"

	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/migrations"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// autoposting-deploy-db will create a database file, with all the tables ready to be used

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

	defer db.Close()
	// Create User table
	migrations.CreateUsers(db)
	// Create Category table
	migrations.CreateCategories(db)
	// Create Posts table
	migrations.CreatePosts(db)

	// Create image and video categories
	image := entities.Category{Name: "image"}
	video := entities.Category{Name: "video"}
	gif := entities.Category{Name: "gif"}
	db.Create(&gif)
	db.Create(&image)
	db.Create(&video)

}
