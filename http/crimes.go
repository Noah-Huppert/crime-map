package http

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/Noah-Huppert/crime-map/models"
)

// QueryParamLimitKey holds the key which the limit query parameter will be
// passed by
const QueryParamLimitKey string = "limit"

// QueryParamOrderByKey holds the key which the order by query parameter will
// be passed by
const QueryParamOrderByKey string = "order_by"

// RespKeyCrimes holds the key which will requests Crime models will be returned
// in
const RespKeyCrimes string = "crimes"

// GetCrimesHandler lists the existing crimes in the database
type GetCrimesHandler struct{}

// Register implements Registerable for GetCrimesHandler
func (h GetCrimesHandler) Register(r *mux.Router) error {
	r.Path("/api/v1/crimes").
		Methods("GET").
		Queries("limit", "{limit:.+}").
		Queries("order_by", "{order_by:.+}").
		Handler(GetCrimesHandler{})

	return nil
}

// ServeHTTP implements the serve method for http.Handler. Requires the request
// contain the 'limit' and 'order_by' query variables.
func (h GetCrimesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get query params
	limit, orderBy, errs := h.parseParams(req)
	if len(errs) != 0 {
		WriteErr(w, errs...)
		return
	}

	// Query
	crimes, err := models.QueryAllCrimes(limit, orderBy)
	if err != nil {
		WriteErr(w, fmt.Errorf("error querying for crimes: %s",
			err.Error()))
		return
	}

	// Response
	resp := make(map[string]interface{})
	resp[RespKeyCrimes] = crimes

	WriteResp(w, resp)
}

// parseParams extracts the 'limit' and 'order_by' query parameters from the
// request. And returns them, along with an array of errors that may have
// occured. This will be len = 0 on success.
func (g GetCrimesHandler) parseParams(req *http.Request) (uint, models.OrderByType, []error) {
	// Record any errors
	errs := []error{}

	// Get vars
	vars := mux.Vars(req)
	var limit uint64 = 0
	var orderBy models.OrderByType = models.OrderByErr

	// If limit query provided
	if query, ok := vars[QueryParamLimitKey]; ok {
		// Convert into uint
		val, err := strconv.ParseUint(query, 10, 64)

		// If error
		if err != nil {
			errs = append(errs, fmt.Errorf("error parsing 'limit' "+
				"query parameter into uint: %s", err.Error()))
		} else {
			// If success
			limit = val
		}
	} else {
		// If not provided
		errs = append(errs, errors.New("'limit' query parameter must "+
			"be provided"))
	}

	// If orderBy query provided
	if query, ok := vars[QueryParamOrderByKey]; ok {
		// Convert to OrderByType
		val, err := models.NewOrderByType(query)

		// If error
		if err != nil {
			errs = append(errs, fmt.Errorf("error parsing 'order_by'"+
				" query parameter into OrderByType: %s",
				err.Error()))
		} else {
			// If success
			orderBy = val
		}
	} else {
		// If not provided
		errs = append(errs, errors.New("'order_by' query parameter must "+
			"be provided"))
	}

	return uint(limit), orderBy, errs
}
