package analysisadapter

import (
	"github.com/shitpostingio/autopostingbot/config/structs"
)

var (
	config structs.AnalysisAPIConfiguration
)

// Start saves a local copy of the Analysis API configuration
// for future use.
func Start(analysisConfig structs.AnalysisAPIConfiguration) {
	config = analysisConfig
}
