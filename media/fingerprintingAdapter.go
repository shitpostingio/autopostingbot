package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/shitposting/autoposting-bot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
	"gitlab.com/shitposting/datalibrary/entities/fpserver"
	"gitlab.com/shitposting/loglog/loglogclient"

	configuration "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

// getPhotoFingerprint asks FPServer the aHash and the pHash of a photo
func getPhotoFingerprint(fileID string, path string, fp configuration.FpServerConfig, botToken string, log *loglogclient.LoglogClient) (aHash string, pHash string, err error) {

	/* SET UP CLIENT AND REQUEST */
	client := &http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", fp.Address, fp.ImageEndpoint, fileID), nil)
	if err != nil {
		log.Err("Can't setup fingerprinting request to FPServer for picture with fileID " + fileID)
		return "", "", err
	}

	/* ADD AUTHENTICATION ERROR AND PERFORM THE REQUEST */
	request.Header.Add(fp.AuthorizationHeaderName, fp.AuthorizationHeaderValue)
	request.Header.Add(fp.CallerAPIKeyHeaderName, botToken)
	request.Header.Add(fp.FilePathHeaderName, path)

	webResponse, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer utility.CloseSafely(webResponse.Body)

	/* READ RESPONSE */
	var response fpserver.Response
	bodyResult, err := ioutil.ReadAll(webResponse.Body)
	if err != nil {
		return "", "", err
	}

	/* UNMARSHAL */
	err = json.Unmarshal(bodyResult, &response)
	return response.AHash, response.PHash, err
}

// getVideoFingerprint asks FPServer the thumbnailFileID, the aHash and the pHash of a video
func getVideoFingerprint(fileInfo *tgbotapi.File, repo *repository.Repository) (thumbnailFileID string, aHash string, pHash string, err error) {

	/* DON'T FINGERPRINT VIDEOS TOO BIG */
	if fileInfo.FileSize > repo.Config.FileSizeThreshold {
		return "", "", "", fmt.Errorf("video size too big")
	}

	// SET UP CLIENT AND REQUEST
	client := &http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", repo.Config.Fpserver.Address, repo.Config.Fpserver.VideoEndpoint, fileInfo.FileID), nil)
	if err != nil {
		repo.Log.Err("Can't setup request to screengen server for picture with fileID " + fileInfo.FileID)
		return "", "", "", err
	}

	// ADD AUTHENTICATION ERROR AND PERFORM THE REQUEST
	request.Header.Add(repo.Config.Fpserver.AuthorizationHeaderName, repo.Config.Fpserver.AuthorizationHeaderValue)
	request.Header.Add(repo.Config.Fpserver.CallerAPIKeyHeaderName, repo.Config.BotToken)
	request.Header.Add(repo.Config.Fpserver.FilePathHeaderName, fileInfo.FilePath)

	webResponse, err := client.Do(request)
	if err != nil {
		return "", "", "", err
	}
	defer utility.CloseSafely(webResponse.Body)

	// READ RESPONSE
	var response fpserver.Response
	bodyResult, err := ioutil.ReadAll(webResponse.Body)
	if err != nil {
		return "", "", "", err
	}

	// UNMARSHAL
	err = json.Unmarshal(bodyResult, &response)
	if err != nil {
		repo.Log.Warn(fmt.Sprintf("Error when unmarshaling FpServer result: %s", string(bodyResult)))
	}

	return response.ThumbnailFileID, response.AHash, response.PHash, err
}

// findPHash returns the matching Fingeprint entity
func findPHash(pHashToFind string, fingerprints []entities.Fingerprint) entities.Fingerprint {

	/* LOOK FOR MATCHING FINGERPRINT IN THE INPUT PARAMETERS */
	for _, fingerprint := range fingerprints {
		if fingerprint.PHash == pHashToFind {
			return fingerprint
		}
	}

	/* RETURN EMPTY FINGERPRINT IF NO MATCH */
	return entities.Fingerprint{}
}
