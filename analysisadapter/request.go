package analysisadapter

import (
	"errors"
	log "github.com/sirupsen/logrus"
	analysis "gitlab.com/shitposting/analysis-api/api/client"
	"gitlab.com/shitposting/analysis-api/services/structs"
	"os"
	"strings"
)

func Request(path, mediaType, fileUniqueID string) (fingerprint *structs.FingerprintResponse, err error) {

	//
	file, err := os.Open(path)
	if err != nil {
		log.Debugln("analisysadapter.Request: unable to open file ", path, ", error: ", err)
		return
	}

	//
	endpoint := getEndpoint(mediaType, fileUniqueID)
	fileNameStart := strings.LastIndex(path, "/") + 1
	filename := path[fileNameStart:]
	log.Debugln("analysisadapter.Request: filename: ", filename, " endpoint: ", endpoint)

	//
	result, errString := analysis.PerformFingerprintRequest(file, filename, endpoint, config.AuthorizationHeaderValue)
	if errString != "" {
		err = errors.New(errString)
	}

	//
	log.Debugln("analysisadapter.Request: result: ", result, " err: ", errString)
	return &result, err

}
