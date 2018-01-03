package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrsKey is the key in HTTP responses which holds the errors array
const ErrsKey string = "errors"

// WriteErr sends an error via a http.ResponseWriter
func WriteErr(w http.ResponseWriter, errs ...error) {
	// Make response
	resp := make(map[string]interface{})
	errsArr := []string{}

	// Add errors
	for _, err := range errs {
		errsArr = append(errsArr, err.Error())
	}
	resp[ErrsKey] = errsArr

	// Marshall into json
	var respStr string = "{\"errors\": []}"
	bytes, err := json.Marshal(resp)
	if err != nil {
		// Manually make json response
		respStr = fmt.Sprintf("{\"errors\": [\"error marshalling errors"+
			"into json: %s\"]", err.Error())
	} else {
		respStr = string(bytes)
	}

	// Send
	fmt.Fprintf(w, respStr)
}
