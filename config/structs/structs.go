package structs

import "go.mongodb.org/mongo-driver/mongo/options"

// Config is a structure containing information used
// to set up the bot
type Config struct {
	BotToken           string
	ChannelID          int64
	FileSizeThreshold  int `type:"optional"`
	MemePath           string
	PostAlertThreshold int `type:"optional"`
	Tdlib              TdlibConfiguration
	DocumentStore      DocumentStoreConfiguration
	AnalysisAPI        AnalysisAPIConfig
	Localization       Localization
}

type AnalysisAPIConfig struct {
	Address                  string
	ImageEndpoint            string
	VideoEndpoint            string
	AuthorizationHeaderName  string
	AuthorizationHeaderValue string
	CallerAPIKeyHeaderName   string
}


