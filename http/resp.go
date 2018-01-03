package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteResp sends successful api response via a http.ResponseWriter
func WriteResp(w http.ResponseWriter, resp map[string]interface{}) {
	// Make response successful
	resp[ErrsKey] = []error{}

	// Marshall into json
	var respStr string = "{\"errors\": []}"
	bytes, err := json.Marshal(resp)
	if err != nil {
		// Manually make json response
		respStr = fmt.Sprintf("{\"errors\": [\"error marshalling "+
			"successful response into json: %s\"]", err.Error())
	} else {
		respStr = string(bytes)
	}

	// Send
	fmt.Fprintf(w, respStr)
}
