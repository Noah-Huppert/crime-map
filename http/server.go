package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/Noah-Huppert/crime-map/config"
)

// Server manges HTTP handlers
type Server struct {
	// router is the Gorilla router used to map requests
	router *mux.Router

	// Routes holds all registered handlers. This field may be manipulated
	// to add and remove handlers. Primarily in the NewServer method.
	Routes []Registerable
}

// NewServer makes a new Server instance
func NewServer() *Server {
	return &Server{
		router: mux.NewRouter(),
		Routes: []Registerable{
			GetCrimesHandler{},
		},
	}
}

// Register adds all handlers to the router. Although simimlar to
// Registerable.Register, it is not the same.
func (s Server) Register() error {
	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// Loop over routes
	for i, route := range s.Routes {
		// Register each route
		if err := route.Register(s.router); err != nil {
			return fmt.Errorf("error registering route i = %d, err"+
				": %s", i, err.Error())
		}
	}

	// Success
	return nil
}

// Serve starts the HTTP server component. An error is returned if one occurs,
// or nil on success
func (s Server) Serve() error {
	// Get config
	c, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %s\n",
			err.Error())
	}

	// Setup routes
	if err = s.Register(); err != nil {
		return fmt.Errorf("error setting up routes: %s", err.Error())
	}

	// Start listening
	fmt.Printf("listening on :%d\n", c.HTTP.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.HTTP.Port), s.router)
}
