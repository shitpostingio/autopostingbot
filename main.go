package main

import (
	"flag"
	"fmt"
	"github.com/bykovme/gotrans"
	"github.com/shitpostingio/autopostingbot/analysisadapter"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/config"
	"github.com/shitpostingio/autopostingbot/documentstore"
	"github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
	"github.com/shitpostingio/autopostingbot/updates"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/autopostingbot/repository"
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

	repository.Config = cfg

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

	// Get the channel chat
	_, err = api.GetChat(cfg.Autoposting.ChannelID)
	if err != nil {
		log.Fatal("Unable to get channel chat")
	}

	// Start listening for updates
	listener := tdlibClient.GetListener()
	defer listener.Close()
	go updates.HandleUpdates(listener)

	// Start the posting manager
	posting.Start(cfg, debug)
	log.Info(fmt.Sprintf("Shitposting autoposting-bot version v%s, build %s, algorithm %s channelname %s", Version, Build, cfg.Autoposting.Algorithm, posting.GetChannelHandle()))
	posting.Listen()

}

func loadCLIParams() {
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.BoolVar(&debug, "debug", false, "activate debug features")
	flag.Parse()
}
