package uniwebhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../models"
)

func HandleWebhook(response http.ResponseWriter, request *http.Request) {

	logger.Debugf("New Request")

	newCall := &models.WebhookCall{
		Timestamp:     time.Now(),
		Host:          request.Host,
		Path:          request.URL.String(),
		Method:        request.Method,
		ContentLength: int(request.ContentLength),
		From:          request.RemoteAddr,
	}

	// Loop through headers
	for name, headers := range request.Header {
		if newCall.Headers == nil {
			newCall.Headers = make(map[string]string)
		}

		var value string
		for i, h := range headers {
			value += h
			if i < (len(headers) - 1) {value += "\n"}
		}
		newCall.Headers[name] = value
	}

	// Loop through queries
	for name, query := range request.URL.Query() {
		if newCall.Query == nil {
			newCall.Query = make(map[string]string)
		}

		var value string
		for i, h := range query {
			value += h
			if i < (len(query) - 1) {value += "\n"}
		}
		newCall.Query[name] = value
	}

	// If this is a POST, add post data
	if request.Method != "GET" {

		err := request.ParseMultipartForm(32 << 20) // 32mb
		if err != nil && err != http.ErrNotMultipart {
			logger.Errorf("Error parsting form: %s", err)
		}

		err = request.ParseForm()
		if err != nil && err != http.ErrNotMultipart {
			logger.Errorf("Error parsting form: %s", err)
		}

		for name, form := range request.PostForm {
			if newCall.PostForm == nil {
				newCall.PostForm = make(map[string]string)
			}

			var value string
			for i, h := range form {
				value += h
				if i < (len(form) - 1) {value += "\n"}
			}
			newCall.PostForm[name] = value
		}

		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			logger.Errorf("Error parsing body: %s", err)
		} else if len(b) > 0 {
			newCall.Body = fmt.Sprintf("%s", b)
		}

		if request.MultipartForm != nil {
			for name, _ := range request.MultipartForm.File {
				if newCall.Files == nil {
					newCall.Files = make(map[string]models.WebhookFile)
				}

				file, handler, err := request.FormFile(name)
				if err != nil {
					logger.Errorf("Error getting file [%s]: %s", name, err)
					break
				}
				defer file.Close()

				newFile := models.WebhookFile{
					Filename: handler.Filename,
					Size:     int(handler.Size),
				}
				newFile.Contents.ReadFrom(file)

				// Loop through headers
				for subName, headers := range handler.Header {
					if newFile.Headers == nil {
						newFile.Headers = make(map[string]string)
					}

					var value string
					for i, h := range headers {
						value += h
						if i < (len(headers) - 1) {value += "\n"}
					}
					newFile.Headers[subName] = value
				}

				newCall.Files[name] = newFile
			}
		}


	}
	go newCall.Save()

	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, "{\"cat\":\"meow\"}")

	return
}

