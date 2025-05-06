package main

import (
	"net/http"

	"github.com/Shriharsh07/chaintrack/config"
	"github.com/Shriharsh07/chaintrack/routes"
	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
