package config

import (
	"github.com/shitpostingio/autopostingbot/config/structs"
	"github.com/spf13/viper"
)

var (
	cfg *structs.Config
	err error
)

// Load reads a configuration file and returns its config instance.
func Load(path string) (*structs.Config, error) {

	//
	setDefaultValuesForOptionalFields()

	//
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	//
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	//
	err = checkMandatoryFields(false, *cfg)
	if err != nil {
		return cfg, err
	}

	//
	err = viper.WriteConfig()
	return cfg, err

}
