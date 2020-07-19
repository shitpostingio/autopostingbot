package structs

// NSFWResponse represents the response of the NSFW service.
type NSFWResponse struct {
	IsNSFW     bool
	Confidence float64
	Label      string
}
