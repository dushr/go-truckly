package truckly

import (
	"encoding/json"

	"google.golang.org/appengine"
)

type Truck struct {
	Name        string   `json:"name"`
	Type        string   `json:"facility_type"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	OpeningTime string   `json:"opening_times"`
	Location    Location `json:"location"`
}

type TruckIndex struct {
	Name        string
	Description string
	Location    appengine.GeoPoint
}

type Location struct {
	Coordinates appengine.GeoPoint `json:"coordinates"`
	Type        string             `json:"type"`
}

type InterimLocation struct {
	Coordinates [2]float64 `json:"coordinates"`
	Type        string     `json:"type"`
}

type Trucks []Truck

type TruckIndexes []TruckIndexes

func (l *Location) UnmarshalJSON(data []byte) error {
	var il InterimLocation

	if err := json.Unmarshal(data, &il); err != nil {
		return err
	}
	l.Coordinates.Lng = il.Coordinates[0]
	l.Coordinates.Lat = il.Coordinates[1]
	l.Type = il.Type
	return nil
}

func (l *Location) MarshalJSON() ([]byte, error) {
	var coor [2]float64
	coor[0] = l.Coordinates.Lng
	coor[1] = l.Coordinates.Lat
	il := InterimLocation{
		Coordinates: coor,
		Type:        l.Type,
	}
	return json.Marshal(&il)
}
