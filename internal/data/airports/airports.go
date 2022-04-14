package airports

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
)

//go:embed airports.json
var airports_json []byte

type AirportEntry struct {
	ICAO string `json:"icao"`
	IATA string `json:"iata"`
	Name string `json:"name"`
}

type airportList struct {
	ByICAO map[string]AirportEntry
	ByIATA map[string]AirportEntry
}

var AirporData *airportList = &airportList{}

func init() {
	dataReader := bufio.NewReader(bytes.NewBuffer(airports_json))
	decoder := json.NewDecoder(dataReader)

	var ICAOList map[string]AirportEntry
	decodeErr := decoder.Decode(&ICAOList)
	if decodeErr != nil {
		panic(decodeErr)
	}

	AirporData.ByICAO = ICAOList
	AirporData.ByIATA = createIATAList(ICAOList)
}

func createIATAList(icaoInput map[string]AirportEntry) map[string]AirportEntry {
	iataList := make(map[string]AirportEntry, len(icaoInput))
	for _, airport := range icaoInput {
		iataList[airport.IATA] = airport
	}
	return iataList
}
