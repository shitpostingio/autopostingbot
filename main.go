package main

import (
	"flag"
	"fmt"
	"github.com/bykovme/gotrans"
	log "github.com/sirupsen/logrus"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/documentstore"
	"gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/posting"
	"gitlab.com/shitposting/autoposting-bot/updates"

	"gitlab.com/shitposting/autoposting-bot/repository"
)

var (
	// configFilePath is the path where the configuration file will be read.
	configFilePath string

	// Version is the autoposting-bot version, a compile-time value.
	Version string

	// Build is the git tag for the current version.
	Build string

	//debug is used to turn on debugging features.
	debug bool

)

func main() {

	// Load parameters from CLI
	loadCLIParams()

	// Set logrus to print debug messages if we're in debug mode
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Load configuration file
	cfg, err := config.Load(configFilePath)
	if err != nil {
		log.Fatal("Error while loading configuration: ", err)
	}

	repository.Config = &cfg

	// Set localization
	err = gotrans.InitLocales(cfg.Localization.Path)
	if err != nil {
		log.Fatal("Error while initializing language files:", err)
	}

	localization.SetLanguage(cfg.Localization.Language)

	// Configure analysis adapter
	analysisadapter.Start(cfg.AnalysisAPI)

	// Connect to the database
	documentstore.Connect(&cfg.DocumentStore, cfg.Autoposting.MediaApproximation, cfg.Autoposting.SimilarityThreshold)

	// Authorize on tdlib
	tdlibClient, err := api.Authorize(cfg.Autoposting.BotToken, &cfg.Tdlib)
	if err != nil {
		log.Fatal("Error while authorizing the bot via tdlib: ", err)
	}

	repository.Tdlib = tdlibClient

	// Get information on self
	repository.Me, err = tdlibClient.GetMe()
	if err != nil {
		log.Fatal("Error while getting information on self from Telegram: ", err)
	}

	// Start listening for updates
	listener := tdlibClient.GetListener()
	defer listener.Close()
	go updates.HandleUpdates(listener)

	// Start the posting manager
	posting.Start(&cfg, debug)
	log.Info(fmt.Sprintf("Shitposting autoposting-bot version v%s, build %s, edition %s", Version, Build, posting.GetPostingManager().GetEditionName()))
	posting.Listen()

}

func loadCLIParams() {
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.BoolVar(&debug, "debug", false, "activate debug features")
	flag.Parse()
}
