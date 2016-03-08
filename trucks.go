package truckly

import (
	"encoding/json"

	"google.golang.org/appengine"
)

type Truck struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"facility_type"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	OpeningTime string   `json:"opening_times"`
	Location    Location `json:"location"`
}

type InterimLocation struct {
	Coordinates appengine.GeoPoint `json:"coordinates"`
	Type        string             `json:"type"`
}

type Location struct {
	Coordinates [2]float64 `json:"coordinates"`
	Type        string     `json:"type"`
}

type Trucks []Truck

func (l *Location) UnmarshalJSON(data []byte) error {
	var il InterimLocation

	if err := json.Unmarshal(data, &il); err != nil {
		return err
	}
	l.Coordinates[0] = il.Coordinates.Lng
	l.Coordinates[1] = il.Coordinates.Lat
	l.Type = il.Type
	return nil
}
