package http

import (
	"github.com/gorilla/mux"
)

type Registerable interface {
	// Register adds a route handler onto a Gorilla mux.Router. An error is
	// returned if one occurs. Or nil on success
	Register(r *mux.Router) error
}
