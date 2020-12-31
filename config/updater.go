package config

import "github.com/spf13/viper"

// UpdateChannelHandle writes the new channel handle to the configuration file
func UpdateChannelHandle(newName string) error {

	//
	cfg.Autoposting.ChannelHandle = newName

	//
	viper.Set("autoposting.channelhandle", newName)
	return viper.WriteConfig()

}
