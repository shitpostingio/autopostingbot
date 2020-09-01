package config

import "github.com/spf13/viper"

const (

	// DocumentStore
	defaultDocumentStoreHosts = "localhost:27017"

	// Tdlib
	defaultTdlibUseTestDc = false
	defaultTdlibUseFileDatabase = true
	defaultTdlibUseChatInfoDatabase = true
	defaultTdlibUseMessageDatabase = true
	defaultTdlibUseSecretChats = false
	defaultTdlibSystemLanguageCode = "en"
	defaultTdlibDeviceModel = "Other"
	defaultTdlibSystemVersion = "1.0.0"
	defaultTdlibApplicationVersion = "1.0.0"
	defaultTdlibEnableStorageOptimizer = true
	defaultTdlibIgnoreFileNames = true

)

func setDefaultValuesForOptionalFields() {

	// DocumentStore
	viper.SetDefault("documentstore.hosts", []string{defaultDocumentStoreHosts})

	// Tdlib
	viper.SetDefault("tdlib.usetestdc", defaultTdlibUseTestDc)
	viper.SetDefault("tdlib.usefiledatabase", defaultTdlibUseFileDatabase)
	viper.SetDefault("tdlib.usechatinfodatabase", defaultTdlibUseChatInfoDatabase)
	viper.SetDefault("tdlib.usemessagedatabase", defaultTdlibUseMessageDatabase)
	viper.SetDefault("tdlib.usesecretchats", defaultTdlibUseSecretChats)
	viper.SetDefault("tdlib.systemlanguagecode", defaultTdlibSystemLanguageCode)
	viper.SetDefault("tdlib.devicemodel", defaultTdlibDeviceModel)
	viper.SetDefault("tdlib.systemversion", defaultTdlibSystemVersion)
	viper.SetDefault("tdlib.applicationversion", defaultTdlibApplicationVersion)
	viper.SetDefault("tdlib.enablestorageoptimizer", defaultTdlibEnableStorageOptimizer)
	viper.SetDefault("tdlib.ignorefilenames", defaultTdlibIgnoreFileNames)

}

