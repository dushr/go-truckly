package truckly

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
	"google.golang.org/appengine/urlfetch"
)

type TruckReturn struct {
	Data Trucks `json:"data"`
}

// guestbookKey returns the key used for all guestbook entries.
func parentTruckKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "TestTruck", "default_truck", 0, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	distance := r.FormValue("distance")
	var query string

	if latitude != "" && longitude != "" && distance != "" {
		query = fmt.Sprintf("distance(Location, geopoint(%s, %s)) < %s", latitude, longitude, distance)
	} else {
		query = ""
	}

	var TruckIds []*datastore.Key
	index, err := search.Open("TestTruck")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for t := index.Search(c, query, nil); ; {
		var ti TruckIndex
		id, err := t.Next(&ti)
		if err == search.Done {
			break
		}
		intID, _ := strconv.Atoi(id)
		key := datastore.NewKey(c, "TestTruck", "", int64(intID), parentTruckKey(c))
		TruckIds = append(TruckIds, key)
	}

	trucks := make([]Truck, len(TruckIds))
	if err := datastore.GetMulti(c, TruckIds, trucks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	truckreturn := TruckReturn{
		Data: trucks,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(truckreturn); err != nil {
		panic(err)
	}
}

func NewTruck(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var truck Truck

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	log.Debugf(c, string(body))
	if err := json.Unmarshal(body, &truck); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	key := datastore.NewIncompleteKey(c, "TestTruck", parentTruckKey(c))
	lol, err := datastore.Put(c, key, &truck)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	index, err := search.Open("TestTruck")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	truckindex := TruckIndex{
		Name:        truck.Name,
		Description: truck.Description,
		Location:    truck.Location.Coordinates,
	}
	_, err = index.Put(c, strconv.FormatInt(lol.IntID(), 10), &truckindex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func ImportTrucks(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get("http://truckly.api.dush.me/trucks/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tr TruckReturn

	if err := json.Unmarshal(body, &tr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, truck := range tr.Data {
		log.Infof(c, truck.Name)
		key := datastore.NewIncompleteKey(c, "TestTruck", parentTruckKey(c))
		lol, err := datastore.Put(c, key, &truck)
		if err != nil {
			log.Errorf(c, err.Error())
		} else {
			log.Infof(c, "Success", lol.StringID())
			index, err := search.Open("TestTruck")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			truckindex := TruckIndex{
				Name:        truck.Name,
				Description: truck.Description,
				Location:    truck.Location.Coordinates,
			}
			_, err = index.Put(c, strconv.FormatInt(lol.IntID(), 10), &truckindex)
			if err != nil {
				log.Errorf(c, err.Error())
			} else {
				log.Infof(c, "Indexed", lol.StringID())
				log.Infof(c, "Indexed", lol.IntID())
			}
		}

	}

	fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
}
