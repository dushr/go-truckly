package truckly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TruckReturn struct {
	Data Trucks `json:"data"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	dat, _ := ioutil.ReadFile("mock/trucks.json")

	var trucks Trucks
	if err := json.Unmarshal(dat, &trucks); err != nil {
		panic(err)
	}
	truckreturn := TruckReturn{
		Data: trucks,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(truckreturn); err != nil {
		panic(err)
	}
}
