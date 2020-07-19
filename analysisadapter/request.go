package analysisadapter

import (
	"errors"
	log "github.com/sirupsen/logrus"
	analysis "gitlab.com/shitposting/analysis-api/api/client"
	"gitlab.com/shitposting/analysis-api/services/structs"
	"os"
	"strings"
)

func Request(path, mediaType, fileUniqueID string) (fingerprint *structs.FingerprintResponse, err error){

	file, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return
	}

	index := strings.LastIndex(path, "/")
	filename := path[index+1:]
	endpoint := getEndpoint(mediaType, fileUniqueID)
	log.Println("Filename: ", filename, " endpoint: ", endpoint)
	result, errString := analysis.PerformFingerprintRequest(file,filename, endpoint, config.AuthorizationHeaderValue)
	if errString != "" {
		err = errors.New(errString)
	}

	log.Println("Risultato:", result, errString)
	return &result, err

}

