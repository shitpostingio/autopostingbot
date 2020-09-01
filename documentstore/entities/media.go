package entities

// Media represents a media in the document store
type Media struct {
	Type             string
	TdlibID          int32
	FileUniqueID     string
	FileID           string
	Histogram        []float64 `bson:",omitempty"`
	HistogramAverage float64   `bson:",omitempty"`
	HistogramSum     float64   `bson:",omitempty"`
	PHash            string    `bson:",omitempty"`
}

// GetHistogramAverageAndSum gets the average and the sum of the input histogram values
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
