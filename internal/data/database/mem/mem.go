package mem

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hschimke/planeTracker/internal/data/model"
)

type MemoryDatabase struct {
	flights []model.Flight
}

func (mdb *MemoryDatabase) GetFlightsForUser(_ context.Context, user model.UserId) ([]model.Flight, error) {
	var returnList []model.Flight
	for _, flight := range mdb.flights {
		if flight.FlightUser == user {
			returnList = append(returnList, flight)
		}
	}
	return returnList, nil
}

func (mdb *MemoryDatabase) AddFlight(_ context.Context, flight model.Flight) (model.FlightId, error) {
	if flight.Id == "" {
		flight.Id = model.FlightId(uuid.New().String())
	}
	mdb.flights = append(mdb.flights, flight)

	return flight.Id, nil
}

func (mdb *MemoryDatabase) DeleteFlight(_ context.Context, flight model.Flight) error {
	index := int(-1)

	for indx, sFlight := range mdb.flights {
		if sFlight.Id == flight.Id {
			index = indx
			break
		}
	}

	if index < 0 {
		return fmt.Errorf("Not Found")
	}

	mdb.flights[index] = mdb.flights[len(mdb.flights)-1]
	mdb.flights = mdb.flights[:len(mdb.flights)-1]
	return nil
}

func (mdb *MemoryDatabase) UpdateFlight(_ context.Context, flight model.Flight) error {
	index := int(-1)

	for indx, sFlight := range mdb.flights {
		if sFlight.Id == flight.Id {
			index = indx
			break
		}
	}

	if index < 0 {
		return fmt.Errorf("Not Found")
	}
	mdb.flights[index] = flight

	return nil
}
