package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
)

const NUM_PARTS = 136

func proxyHandler(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		resp, err := http.Get(target)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, string(body))
	}
}

func registerEndpoint(response http.ResponseWriter, request *http.Request) {
	var endpoint Endpoint
	if request != nil && request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &endpoint); err != nil {
			response.WriteHeader(http.StatusBadRequest)
		}
		http.HandleFunc(endpoint.Origin, proxyHandler(endpoint.Target))
		response.WriteHeader(http.StatusCreated)

	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	bodies := []string{}
	for i := 0; i < NUM_PARTS; i++ {
		resp, err := http.Get(fmt.Sprintf("/part/%d", i))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		bodies = append(bodies, string(body))
	}
	fmt.Fprintf(w, strings.Join(bodies, ""))
}

func main() {

	for i := 0; i < NUM_PARTS; i++ {
		http.HandleFunc(fmt.Sprintf("/part/%d", i), notImplemented)
	}

	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/register_endpoint", registerEndpoint)
	http.ListenAndServe(":8080", nil)
}
