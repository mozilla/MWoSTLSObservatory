package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *CORSRouter{

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return &CORSRouter{router}
}

type CORSRouter struct {
  r *mux.Router
}


func (s *CORSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Content-Type", "application/json")

  if r.Method == "OPTIONS" {
    return
  }

  s.r.ServeHTTP(w, r)
}


type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Scan",
		"POST",
		"/api/v1/scan",
		ScanHandler,
	},
	Route{
		"Results",
		"GET",
		"/api/v1/results",
		ResultHandler,
	},
	Route{
		"Certificate",
		"GET",
		"/api/v1/certificate",
		CertificateHandler,
	},
}
