package client

import (
	"bytes"
	"encoding/json"
	"gitlab.com/shitposting/analysis-api/services/structs"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// PerformAnalysisRequest performs a request to the fingerprinting service.
func PerformAnalysisRequest(file io.Reader, fileName, endpoint, authorization string) (data structs.Analysis, errorString string) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		errorString = err.Error()
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		errorString = err.Error()
		return
	}

	err = writer.Close()
	if err != nil {
		errorString = err.Error()
		return
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		errorString = err.Error()
		return
	}

	// We want to send data with the multipart/form-data Content-Type
	//TODO: sistemare autorizzazione
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-shitposting-key", authorization)
	client := http.Client{Timeout: time.Second * 30}
	response, err := client.Do(request)
	if err != nil {
		errorString = err.Error()
		return
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Println("Analysis.PerformAnalysisRequest: unable to close response body", err)
		}
	}()

	bodyResult, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorString = err.Error()
		return
	}

	var ar structs.Analysis
	err = json.Unmarshal(bodyResult, &ar)
	if err != nil {
		errorString = err.Error()
		log.Println("PerformAnalysisRequest: error while unmarshaling ", err)
		return
	}

	if ar.FingerprintErrorString != "" {
		errorString = ar.FingerprintErrorString
	} else {
		errorString = ar.NSFWErrorString
	}

	return ar, errorString

}
