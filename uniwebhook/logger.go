package uniwebhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../models"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		newCall := &models.WebhookCall{
			Timestamp:     time.Now(),
			Host:          r.Host,
			Path:          r.URL.String(),
			Method:        r.Method,
			ContentLength: int(r.ContentLength),
			From:          r.RemoteAddr,
		}

		// Loop through headers
		for name, headers := range r.Header {
			if newCall.Headers == nil {
				newCall.Headers = make(map[string]string)
			}

			var value string
			for i, h := range headers {
				value += h
				if i < (len(headers) - 1) {
					value += "\n"
				}
			}
			newCall.Headers[name] = value
		}

		// Loop through queries
		for name, query := range r.URL.Query() {
			if newCall.Query == nil {
				newCall.Query = make(map[string]string)
			}

			var value string
			for i, h := range query {
				value += h
				if i < (len(query) - 1) {
					value += "\n"
				}
			}
			newCall.Query[name] = value
		}

		// If this is a POST, add post data
		if r.Method != "GET" {

			err := r.ParseMultipartForm(32 << 20) // 32mb
			if err != nil && err != http.ErrNotMultipart {
				logger.Errorf("Error parsting form: %s", err)
			}

			err = r.ParseForm()
			if err != nil && err != http.ErrNotMultipart {
				logger.Errorf("Error parsting form: %s", err)
			}

			for name, form := range r.PostForm {
				if newCall.PostForm == nil {
					newCall.PostForm = make(map[string]string)
				}

				var value string
				for i, h := range form {
					value += h
					if i < (len(form) - 1) {
						value += "\n"
					}
				}
				newCall.PostForm[name] = value
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Errorf("Error parsing body: %s", err)
			} else if len(b) > 0 {
				newCall.Body = fmt.Sprintf("%s", b)
			}

			if r.MultipartForm != nil {
				for name, _ := range r.MultipartForm.File {
					if newCall.Files == nil {
						newCall.Files = make(map[string]models.WebhookFile)
					}

					file, handler, err := r.FormFile(name)
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
							if i < (len(headers) - 1) {
								value += "\n"
							}
						}
						newFile.Headers[subName] = value
					}

					newCall.Files[name] = newFile
				}
			}

		}
		go newCall.Save()

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}