package structs

// AutopostingConfiguration represents the configuration of the autoposting bot.
type AutopostingConfiguration struct {

	// BotToken is the Telegram bot token we will need to log in to.
	BotToken string

	// ChannelID is the id of the channel on which we will send posts.
	ChannelID int64

	// FileSizeThreshold represents the maximum file size we can perform
	// fingerprint requests on.
	FileSizeThreshold int `type:"optional"`

	// MediaPath is the path where we will save posted medias.
	MediaPath string

	// PostAlertThreshold represents the threshold below which the admins
	// will be notified of low posts enqueued.
	PostAlertThreshold int `type:"optional"`

	// MediaApproximation represents the approximation the similarity search
	// in the database will use.
	MediaApproximation float64 `type:"optional"`

	// SimilarityThreshold represents the similarity threshold that tells us
	// whether two pictures are similar enough or not.
	SimilarityThreshold int `type:"optional"`

	// Algorithm represents the algorithm that will be run.
	Algorithm string

	// ChannelHandle represents the tag for the channel
	ChannelHandle string
}
