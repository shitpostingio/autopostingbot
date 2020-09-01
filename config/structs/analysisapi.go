package structs

type AnalysisAPIConfiguration struct {
	Address                  string
	ImageEndpoint            string
	VideoEndpoint            string
	AuthorizationHeaderName  string
	AuthorizationHeaderValue string
	CallerAPIKeyHeaderName   string
}
