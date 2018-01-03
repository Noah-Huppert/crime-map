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
	// Get vars
	vars := mux.Vars(req)
	var limit uint64 = 0
	var orderBy models.OrderByType = models.OrderByErr

	// Parse limit
	if val, ok := vars[QueryParamLimitKey]; ok {
		conv, err := strconv.ParseUint(val, 10, 64)
		limit = conv
		if err != nil {
			WriteErr(w, fmt.Errorf("error parsing 'limit' "+
				"query parameter into uint: %s", err.Error()))
			return
		}
	} else {
		WriteErr(w, errors.New("'limit' query parameter must "+
			"be provided"))
		return
	}

	// Parse orderBy
	if val, ok := vars[QueryParamOrderByKey]; ok {
		conv, err := models.NewOrderByType(val)
		orderBy = conv
		if err != nil {
			WriteErr(w, fmt.Errorf("error parsing 'order_by'"+
				" query parameter into OrderByType: %s",
				err.Error()))
			return
		}
	} else {
		WriteErr(w, errors.New("'order_by' query parameter must "+
			"be provided"))
		return
	}

	// Query
	crimes, err := models.QueryAllCrimes(uint(limit), orderBy)
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
