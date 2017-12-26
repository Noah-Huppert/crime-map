package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Noah-Huppert/crime-map/config"
	"github.com/Noah-Huppert/crime-map/dstore"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/parsers"

	"github.com/gorilla/mux"
)

const file = "data/2017-10-12.pdf"

func main() {
	// Get config
	fmt.Println("loading configuration")
	c, err := config.NewConfig()
	if err != nil {
		fmt.Printf("error loading configuration: %s\n", err.Error())
		os.Exit(1)
	}

	// Migrate db
	fmt.Println("migrating db")
	err = models.Migrate()
	if err != nil {
		fmt.Printf("error migrating db: %s\n", err.Error())
		os.Exit(1)
	}

	// Parse crimes
	fmt.Println("parsing report")
	r := parsers.NewReport(file)

	crimes, err := r.Parse()
	if err != nil {
		fmt.Printf("error parsing report: %s\n", err.Error())
		os.Exit(1)
	}

	// Print crimes
	for _, crime := range crimes {
		if err = crime.SaveIfUnique(&crime, &crime); err != nil {
			fmt.Printf("error saving crime: %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Start http server
	router := mux.NewRouter()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	router.HandleFunc("/api/v1/crimes", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(crimes)
		if err != nil {
			fmt.Fprintf(w, "error marshalling crimes: %s", err.Error())
		}

		fmt.Fprintf(w, string(bytes))
	})

	fmt.Printf("listening on :%d\n", c.HTTP.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.HTTP.Port), router)
	if err != nil {
		fmt.Printf("error starting http server: %s\n", err.Error())
	}
}
