package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hschimke/planeTracker/internal/data/database/mem"
	"github.com/hschimke/planeTracker/internal/routes"
)

const apiString = "/api/v1/"

func main() {
	//dbString := os.Getenv("CONNECTION_STRING")
	serverPort := os.Getenv("SERVER_PORT")

	//flightDatabase := pg.NewPostgresDatabase(dbString)
	flightDatabase := &mem.MemoryDatabase{}
	routerMap := routes.NewServer(flightDatabase)

	router := http.NewServeMux()
	router.Handle(apiString+"getFlights", authRequiredMW(http.HandlerFunc(routerMap.GetFlightsForUser)))
	router.Handle(apiString+"addFlight", authRequiredMW(http.HandlerFunc(routerMap.AddFlight)))
	router.Handle(apiString+"deleteFlight", authRequiredMW(http.HandlerFunc(routerMap.DeleteFlight)))
	router.Handle(apiString+"updateFlight", authRequiredMW(http.HandlerFunc(routerMap.UpdateFlight)))

	address := fmt.Sprintf(":%s", serverPort)

	server := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	log.Fatal(server.ListenAndServe())
}
