package fpserver

//Response is what we reply to a client that asks us if a media is a duplicate.
type Response struct {
	AHash           string    `json:"ahash"`
	PHash           string    `json:"phash"`
	ThumbnailFileID string    `json:"thumbnail_file_id,omitempty"`
	Histogram       []float64 `json:"histogram"`
}
