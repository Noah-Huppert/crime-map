package http

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/Noah-Huppert/crime-map/models"
)

// QueryParamOffsetKey holds the key which the offset query parameter will be
// passed by
const QueryParamOffsetKey string = "offset"

// QueryParamLimitKey holds the key which the limit query parameter will be
// passed by
const QueryParamLimitKey string = "limit"

// QueryParamOrderByKey holds the key which the order by query parameter will
// be passed by
const QueryParamOrderByKey string = "order_by"

// RespKeyCrimes holds the key which will requests Crime models will be returned
// in
const RespKeyCrimes string = "crimes"

// GetCrimesHandler lists the existing crimes in the database. It expects
// the following query parameters:
//
//	- offset (uint): Index of first element to return, 0 would return the
//			 the first item, 10 would return the 10th item.
// 	- limit (uint): Index of last element to return.
//	- order_by (date_occurred|date_reported): Specifies how to order
//					          returned results.
type GetCrimesHandler struct{}

// Register implements Registerable for GetCrimesHandler
func (h GetCrimesHandler) Register(r *mux.Router) error {
	r.Path("/api/v1/crimes").
		Methods("GET").
		Queries(QueryParamOffsetKey, fmt.Sprintf("{%s:.+}",
			QueryParamOffsetKey)).
		Queries(QueryParamLimitKey, fmt.Sprintf("{%s:.+}",
			QueryParamLimitKey)).
		Queries(QueryParamOrderByKey, fmt.Sprintf("{%s:.+}",
			QueryParamOrderByKey)).
		Handler(GetCrimesHandler{})

	return nil
}

// ServeHTTP implements the serve method for http.Handler. Requires the request
// contain the 'limit' and 'order_by' query variables.
func (h GetCrimesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get query params
	offset, limit, orderBy, errs := h.parseParams(req)
	if len(errs) != 0 {
		WriteErr(w, errs...)
		return
	}

	// Check offset < limit
	if offset >= limit {
		WriteErr(w, fmt.Errorf("'offset' query parameter must be "+
			"less than 'limit' query parameter"))
		return
	}

	// Query
	crimes, err := models.QueryAllCrimes(offset, limit, orderBy)
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

// parseParams extracts the 'offset', 'limit' and 'order_by' query parameters
// from the request. And returns them, along with an array of errors that may
// have occurred. This will be len = 0 on success.
//
// Values returned in the following order: offset, limit, order_by
func (g GetCrimesHandler) parseParams(req *http.Request) (uint, uint, models.OrderByType, []error) {
	// Record any errors
	errs := []error{}

	// Get vars
	vars := mux.Vars(req)
	var offset uint64 = 0
	var limit uint64 = 0
	var orderBy models.OrderByType = models.OrderByErr

	// If offset query provided
	if query, ok := vars[QueryParamOffsetKey]; ok {
		// Convert into uint
		val, err := strconv.ParseUint(query, 10, 64)

		// If error
		if err != nil {
			errs = append(errs, fmt.Errorf("error parsing 'offset' "+
				"query parameter into uint: %s", err.Error()))
		} else {
			// If success
			limit = val
		}
	} else {
		// If not provided
		errs = append(errs, errors.New("'offset' query parameter must "+
			"be provided"))
	}

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

	return uint(offset), uint(limit), orderBy, errs
}
