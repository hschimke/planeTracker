package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hschimke/planeTracker/internal/data/database/pg"
	"github.com/hschimke/planeTracker/internal/routes"
)

const apiString = "/api/v1/"

func main() {
	dbString := os.Getenv("CONNECTION_STRING")
	serverPort := os.Getenv("SERVER_PORT")

	flightDatabase := pg.NewPostgresDatabase(dbString)
	//flightDatabase := &mem.MemoryDatabase{}
	routerMap := routes.NewServer(flightDatabase)

	router := http.NewServeMux()

	// Flight features
	router.Handle(apiString+"getFlights", authRequiredMW(http.HandlerFunc(routerMap.GetFlightsForUser)))
	router.Handle(apiString+"addFlight", authRequiredMW(http.HandlerFunc(routerMap.AddFlight)))
	router.Handle(apiString+"deleteFlight", authRequiredMW(http.HandlerFunc(routerMap.DeleteFlight)))
	router.Handle(apiString+"updateFlight", authRequiredMW(http.HandlerFunc(routerMap.UpdateFlight)))
	router.Handle(apiString+"bulkAddFlights", authRequiredMW(http.HandlerFunc(routerMap.BulkAddFlights)))
	router.Handle(apiString+"getPlaneDetails", authRequiredMW(http.HandlerFunc(routerMap.GetPlaneDetail)))

	// Passenger Features
	router.Handle(apiString+"getPassengersForUser", authRequiredMW(http.HandlerFunc(routerMap.GetPassengersForUser)))
	router.Handle(apiString+"addPassengerForUser", authRequiredMW(http.HandlerFunc(routerMap.AddPassengerForUser)))
	router.Handle(apiString+"removePassengerForUser", authRequiredMW(http.HandlerFunc(routerMap.RemovePassengerForUser)))
	router.Handle(apiString+"setDefaultStatusForPassengerOfUser", authRequiredMW(http.HandlerFunc(routerMap.SetDefaultStatusForPassengerOfUser)))
	router.Handle(apiString+"addPassengerToFlight", authRequiredMW(http.HandlerFunc(routerMap.AddPassengerToFlight)))
	router.Handle(apiString+"removePassengerFromFlight", authRequiredMW(http.HandlerFunc(routerMap.RemovePassengerFromFlight)))
	router.Handle(apiString+"getFlightsAsPassenger", authRequiredMW(http.HandlerFunc(routerMap.GetFlightsAsPassenger)))

	address := fmt.Sprintf(":%s", serverPort)

	server := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	log.Fatal(server.ListenAndServe())
}
