package fpserver

// NSFWResponse represents what fpserver returns after checking if a media is nsfw or not.
type NSFWResponse struct {
	IsNSFW      bool    `json:"is_nsfw"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}
