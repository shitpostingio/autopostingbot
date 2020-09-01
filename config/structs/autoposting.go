package structs

type AutopostingConfiguration struct {
	BotToken            string
	ChannelID           int64
	FileSizeThreshold   int `type:"optional"`
	MediaPath           string
	PostAlertThreshold  int     `type:"optional"`
	MediaApproximation  float64 `type:"optional"`
	SimilarityThreshold int     `type:"optional"`
}
