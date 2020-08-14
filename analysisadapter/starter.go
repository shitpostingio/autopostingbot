package analysisadapter

import (
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	config structs.AnalysisAPIConfig
)

func Start(analysisConfig structs.AnalysisAPIConfig) {
	config = analysisConfig
}
