package http

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/Noah-Huppert/crime-map/models"
)

// ReportsKey is the key which the list of reports will be returned in the
// endpoint
const ReportsKey string = "reports"

// ListReportsHandler retrieves all report models
type ListReportsHandler struct{}

// Register implements the Registerable interface for ListReportsHandler
func (h ListReportsHandler) Register(r *mux.Router) error {
	r.Path("/api/v1/reports").Handler(ListReportsHandler{})

	return nil
}

// ServeHTTP returns a list of Report models in the 'reports' field.
func (h ListReportsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Query
	reports, err := models.QueryAllReports()
	if err != nil {
		WriteErr(w, err)
		return
	}

	// Respond
	resp := make(map[string]interface{})
	resp[ReportsKey] = reports

	WriteResp(w, resp)
	return
}
