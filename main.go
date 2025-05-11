package main

import (
	"net/http"

	"github.com/Shriharsh07/chaintrack/config"
	_ "github.com/Shriharsh07/chaintrack/docs"
	"github.com/Shriharsh07/chaintrack/routes"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	http.ListenAndServe(":8080", r)
}
