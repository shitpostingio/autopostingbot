package config

import (
	"gitlab.com/shitposting/autoposting-bot/config/structs"
	"log"

	"github.com/spf13/viper"

)

const (
	defaultFileSizeThreshold  = 20971520 //20MB
	defaultDatabaseAddress    = "localhost"
	defaultDatabasePort       = 3306
	defaultDocumentStoreHosts = "localhost:27017"
	defaultSocketPath         = "/tmp/loglog.socket"
)

// Load reads a configuration file and returns its config instance
func Load(path string) (cfg structs.Config, err error) {

	setDefaultValuesForOptionalFields()

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	err = CheckMandatoryFields(false, cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = viper.WriteConfig()
	return
}

func setDefaultValuesForOptionalFields() {
	viper.SetDefault("fpserver.filesizethreshold", defaultFileSizeThreshold)
	viper.SetDefault("loglog.socketpath", defaultSocketPath)
	viper.SetDefault("database.address", defaultDatabaseAddress)
	viper.SetDefault("database.port", defaultDatabasePort)
	viper.SetDefault("documentstore.hosts", []string{defaultDocumentStoreHosts})
}
