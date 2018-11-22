package web

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

var codeTitle = map[int]string{
	1:    "Malformed JSON Body",
	2201: "Missing Required Attribute",
	2202: "Requested Relationship Not Found",
}

func HandleNotFound(response http.ResponseWriter, request *http.Request) {
	MakeErrorResponse(response, http.StatusNotFound, request.URL.Path, 0)
	return
}

func MakeErrorResponse(response http.ResponseWriter, status int, detail string, code int) {
	// Make Strings
	var statusStr string = strconv.Itoa(status)
	var codeStr string

	// Get Title
	var title string
	if code == 0 { // code 0 means no code
		title = http.StatusText(status)
	} else {
		title = codeTitle[code]
		codeStr = strconv.Itoa(code)
	}

	// Send Response
	response.WriteHeader(status)

	response.Header().Set("Content-Type", jsonapi.MediaType)
	jsonapi.MarshalErrors(response, []*jsonapi.ErrorObject{{
		Title:  title,
		Detail: detail,
		Status: statusStr,
		Code:   codeStr,
	}})

	return
}

func HandleNotImplemented(response http.ResponseWriter, request *http.Request) {
	MakeErrorResponse(response, http.StatusMethodNotAllowed, request.Method, 0)
	return
}