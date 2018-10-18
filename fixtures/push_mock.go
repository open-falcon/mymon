// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// CreateEndpoint write data of mymon pushing to push_metric.txt
func CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	metrics, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("./push_metric.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(metrics); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/push", CreateEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":1988", router))
}
