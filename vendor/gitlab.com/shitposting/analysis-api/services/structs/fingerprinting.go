package structs

// FingerprintResponse represents the response of the fingerprinting service.
type FingerprintResponse struct {
	PHash     string
	Histogram []float64
}
