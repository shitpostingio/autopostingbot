package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

// WatchConfig monitors the configuration for changes.
func WatchConfig(cfg *structs.Config) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {

		// we need to make sure the config is valid
		// use a temporary one to check it
		var tempCfg structs.Config
		err := viper.Unmarshal(&tempCfg)
		if err != nil {
			log.Error("The configuration file was changed but it couldn't be unmarshalled")
			return
		}

		//Oddly, for each config change, two events seem to be triggered.
		//The first one will have an empty configuration, causing an error here.
		err = checkMandatoryFields(true, tempCfg)
		if err != nil {
			log.Error("The configuration file was changed but there were issues:", err)
			return
		}

		//The config is correct
		err = viper.Unmarshal(&cfg)
		if err != nil {
			panic("The new configuration couldn't be set")
		}

		log.Info("The configuration was updated correctly")

	})
}
