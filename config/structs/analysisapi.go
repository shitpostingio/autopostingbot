package structs

// AnalysisAPIConfiguration represents the configuration of the Analysis API adapter.
type AnalysisAPIConfiguration struct {

	// Address is the http(s) address of the API.
	Address                  string

	// ImageEndpoint is the API endpoint for images.
	// It must not contain the Address part.
	ImageEndpoint            string

	// VideoEndpoint is the API endpoint for videos.
	// It must not contain the Address part.
	VideoEndpoint            string

	// AuthorizationHeaderName is the name of the authorization header
	// that will be checked by the API endpoint.
	AuthorizationHeaderName  string

	// AuthorizationHeaderValue is the authorization token to add in the
	// AuthorizationHeaderName.
	AuthorizationHeaderValue string

	// CallerAPIKeyHeaderName is the Telegram Bot Token of the caller.
	CallerAPIKeyHeaderName   string

}
