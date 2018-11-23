package uniwebhook

import (
	"fmt"
	"net/http"
)

func HandleCat(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, "{\"cat\":\"meow\"}")

	return
}
