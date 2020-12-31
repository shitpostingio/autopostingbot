package config

import "github.com/spf13/viper"

const (

	// Autoposting
	defaultAutopostingFileSizeThreshold   = 20 << 10
	defaultAutopostingPostAlertThreshold  = 10
	defaultAutopostingMediaApproximation  = 0.08
	defaultAutopostingSimilarityThreshold = 6
	defaultAutopostingChannelHandle       = "shitpost"

	// DocumentStore
	defaultDocumentStoreHosts             = "localhost:27017"
	defaultDocumentStoreAuthMechanism     = "SCRAM-SHA-1"
	defaultDocumentStoreUseAuthentication = false
	defaultDocumentStoreUseReplicaSet     = false

	// Tdlib
	defaultTdlibUseTestDc              = false
	defaultTdlibUseFileDatabase        = true
	defaultTdlibUseChatInfoDatabase    = true
	defaultTdlibUseMessageDatabase     = true
	defaultTdlibUseSecretChats         = false
	defaultTdlibSystemLanguageCode     = "en"
	defaultTdlibDeviceModel            = "Other"
	defaultTdlibSystemVersion          = "1.0.0"
	defaultTdlibApplicationVersion     = "1.0.0"
	defaultTdlibEnableStorageOptimizer = true
	defaultTdlibIgnoreFileNames        = true
)

// setDefaultValuesForOptionalFields sets default values for configuration files.
// These values will be persisted in the configuration file after validation.
func setDefaultValuesForOptionalFields() {

	// Autoposting
	viper.SetDefault("autoposting.filesizethreshold", defaultAutopostingFileSizeThreshold)
	viper.SetDefault("autoposting.postalertthreshold", defaultAutopostingPostAlertThreshold)
	viper.SetDefault("autoposting.mediaapproximation", defaultAutopostingMediaApproximation)
	viper.SetDefault("autoposting.similaritythreshold", defaultAutopostingSimilarityThreshold)
	viper.SetDefault("autoposting.channelhandle", defaultAutopostingChannelHandle)

	// DocumentStore
	viper.SetDefault("documentstore.hosts", []string{defaultDocumentStoreHosts})
	viper.SetDefault("documentstore.useauthentication", defaultDocumentStoreUseAuthentication)
	viper.SetDefault("documentstore.usereplicaset", defaultDocumentStoreUseReplicaSet)
	viper.SetDefault("documentstore.authmechanism", defaultDocumentStoreAuthMechanism)

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
