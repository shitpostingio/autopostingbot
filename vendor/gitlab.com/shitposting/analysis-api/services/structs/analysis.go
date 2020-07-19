package structs

// Analysis represents the fingerprint and NSFW data about a media.
type Analysis struct {
	Fingerprint            FingerprintResponse
	NSFW                   NSFWResponse
	FingerprintErrorString string
	NSFWErrorString        string
}
