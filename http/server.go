package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/Noah-Huppert/crime-map/config"
)

// Serve starts the HTTP server component. An error is returned if one occurs,
// or nil on success
func Serve() error {
	// Get config
	fmt.Println("loading configuration")
	c, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %s\n",
			err.Error())
	}

	// Setup routes
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/crimes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "CRIMES")
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	fmt.Printf("listening on :%d\n", c.HTTP.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.HTTP.Port), router)
}
