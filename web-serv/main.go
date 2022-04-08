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
	router.Handle(apiString+"getFlights", authRequiredMW(http.HandlerFunc(routerMap.GetFlightsForUser)))
	router.Handle(apiString+"addFlight", authRequiredMW(http.HandlerFunc(routerMap.AddFlight)))
	router.Handle(apiString+"deleteFlight", authRequiredMW(http.HandlerFunc(routerMap.DeleteFlight)))
	router.Handle(apiString+"updateFlight", authRequiredMW(http.HandlerFunc(routerMap.UpdateFlight)))

	server := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	server.ListenAndServe()
}
