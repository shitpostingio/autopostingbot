package entities

// Media represents a media in the document store.
type Media struct {

	// Type is the media type.
	Type             string

	// TdlibID is the id of the local tdlib database.
	// It is required to download files but it is not very meaningful.
	TdlibID          int32

	// FileUniqueID is Telegram's unique ID of the media.
	FileUniqueID     string

	// FileID is Telegram's ID of the media.
	// This ID is just temporary.
	FileID           string

	// Histogram is an approximate color histogram of the content of the media.
	Histogram        []float64 `bson:",omitempty"`

	// HistogramAverage is Histogram's average.
	HistogramAverage float64   `bson:",omitempty"`

	// HistogramSum is Histogram's weighed sum.
	HistogramSum     float64   `bson:",omitempty"`

	// PHash is the media's perception hash.
	PHash            string    `bson:",omitempty"`
}

// GetHistogramAverageAndSum gets the average and the sum of the input histogram values.
func GetHistogramAverageAndSum(histogram []float64) (average, sum float64) {

	coefficient := 1.0
	for i := 0; i < 16; i++ {
		sum += histogram[i] * coefficient
		sum += histogram[31-i] * coefficient
		coefficient++
	}

	average = sum / float64(len(histogram))
	return

}
