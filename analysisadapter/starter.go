package analysisadapter

import (
	configuration "gitlab.com/shitposting/autoposting-bot/config"
)

var (
	config configuration.AnalysisAPIConfig
)

func Start(analysisConfig configuration.AnalysisAPIConfig) {
	config = analysisConfig
}
