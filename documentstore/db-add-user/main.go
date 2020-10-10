package main

import (
	"flag"
	"github.com/shitpostingio/autopostingbot/config"
	"github.com/shitpostingio/autopostingbot/documentstore"
	"github.com/shitpostingio/autopostingbot/repository"
	log "github.com/sirupsen/logrus"
)

var (
	// config file path, if not specified it will read ./config.toml
	configFilePath string

	//
	userID int
)

func main() {

	// Load parameters from CLI
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.IntVar(&userID, "userid", 0, "User ID to add to the admins")
	flag.Parse()

	//
	if userID == 0 {
		log.Fatal("User ID 0")
	}

	// Load configuration file
	cfg, err := config.Load(configFilePath)
	if err != nil {
		log.Fatal("Error while loading configuration: ", err)
	}

	repository.Config = &cfg

	// Connect to the database
	documentstore.Connect(&cfg.DocumentStore, cfg.Autoposting.MediaApproximation, cfg.Autoposting.SimilarityThreshold)

	//
	err = documentstore.AddUser(int32(userID), documentstore.UserCollection)
	if err != nil {
		log.Fatal("Unable to add user to the database: ", err)
	}

	log.Println("User added correctly!")

}
