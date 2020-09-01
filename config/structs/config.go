package structs

// Config is a structure containing information used
// to set up the bot
type Config struct {
	Autoposting   AutopostingConfiguration
	Tdlib         TdlibConfiguration
	DocumentStore DocumentStoreConfiguration
	AnalysisAPI   AnalysisAPIConfiguration
	Localization  LocalizationConfiguration
}
