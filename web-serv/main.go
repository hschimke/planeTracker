package main

import (
	"github.com/hschimke/planeTracker/internal/data/database/pg"
	"github.com/hschimke/planeTracker/internal/routes"
	"net/http"
	"os"
	"time"
)

const apiString = "/api/v1/"

func main() {
	dbString := os.Getenv("CONNECTION_STRING")
	address := os.Getenv("SERVER_PORT")

	flightDatabase := pg.NewPostgresDatabase(dbString)
	routerMap := routes.NewServer(flightDatabase)

	router := http.NewServeMux()
	router.HandleFunc(apiString+"getFlights", routerMap.GetFlightsForUser)
	router.HandleFunc(apiString+"addFlight", routerMap.AddFlight)
	router.HandleFunc(apiString+"deleteFlight", routerMap.DeleteFlight)
	router.HandleFunc(apiString+"updateFlight", routerMap.UpdateFlight)

	server := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	server.ListenAndServe()
}
