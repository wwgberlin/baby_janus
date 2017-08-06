package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
)

func genericHandler(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		resp, err := http.Get(target)
		if err != nil {
			// handle error
			//status 500
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
		} else {
			if err := json.Unmarshal(body, &endpoint); err != nil {
				response.WriteHeader(http.StatusBadRequest)
			} else {
				http.HandleFunc(endpoint.Origin, genericHandler(endpoint.Target))
				response.WriteHeader(http.StatusCreated)
			}
		}
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 135; i++{
		http.Get(fmt.Sprintf("/part/%d", i))
	}
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/register_endpoint", registerEndpoint)
	http.ListenAndServe(":8080", nil)
}
