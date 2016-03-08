package truckly

import "net/http"

func init() {
	router := NewRouter()
	http.Handle("/", router)
}
