package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"

	"github.com/wwgberlin/baby_janus/api/cluster"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
)

/*
	redirectHandler returns a handler to redirect the request to
 */

func redirectHandler(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusFound)
	}
}

/*
	registerEndpoint handles requests to registers routes origin - target
 */
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
		http.HandleFunc(endpoint.Origin, redirectHandler(endpoint.Target))
		response.WriteHeader(http.StatusCreated)

	}
}

/*
	helloUser fetches the parts from all the APIs registered to the cluster
 */
func helloUser(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodies := []string{}
		for _, path := range c.GetSlices() {
			resp, err := http.Get(path)
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
}

/*
	incrClusterId - returns handler to increment the cluster servers size
 */

func incrClusterId(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.IncrClusterId()))
	}
}

func getInstanceSlices(c cluster.Cluster) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err = strconv.Atoi(string(body));
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slicesJson, _ := json.Marshal(c.GetInstanceSlices(id))
		w.Write(slicesJson)
	}
}

func main() {
	c := cluster.NewCluster()

	/*
		register your initial routes for the API here
	 */
	http.HandleFunc("/", helloUser(c))
	http.HandleFunc("/next_cluster_id", incrClusterId(c))
	http.HandleFunc("/get_instance_slices", getInstanceSlices(c))
	http.HandleFunc("/register_endpoint", registerEndpoint)

	http.ListenAndServe(":8080", nil)

}
