package structs

// Config is the structure that contains all sub-configurations.
type Config struct {

	// Autoposting contains the bot configuration.
	Autoposting   AutopostingConfiguration

	// Tdlib contains tdlib-specific configuration values.
	Tdlib         TdlibConfiguration

	// DocumentStore contains MongoDB configuration values.
	DocumentStore DocumentStoreConfiguration

	// AnalysisAPI contains Analysis API configuration values.
	AnalysisAPI   AnalysisAPIConfiguration

	// Localization contains localization configuration values.
	Localization  LocalizationConfiguration
}
