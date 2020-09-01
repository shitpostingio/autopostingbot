package analysisadapter

import (
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	config structs.AnalysisAPIConfiguration
)

func Start(analysisConfig structs.AnalysisAPIConfiguration) {
	config = analysisConfig
}
