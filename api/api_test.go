package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"runtime"
)

type (
	clusterMock struct {
		iterator func() int
	}
)

func TestRedirect(t *testing.T) {
	startServer()

	calledBack := false
	origin := "/origin"
	target := "http://127.0.0.1:8080/target"

	http.HandleFunc("/target", func(w http.ResponseWriter, r *http.Request) {
		calledBack = true
	})

	if b, err := json.Marshal(struct {
		Origin string;
		Target string
	}{Origin: origin, Target: target}); err != nil {
		t.Fatal(err.Error())
	} else if _, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/register_endpoint"), "application/json", bytes.NewBuffer(b)); err != nil {
		t.Fatal(err.Error())
	}

	if _, err := http.Get("http://127.0.0.1:8080/origin"); err != nil {
		t.Fatal(err.Error())
	}
	if !calledBack {
		t.Error("didn't redirect to target")
	}
}

func TestGetInstanceSlices(t *testing.T) {
	mock := clusterMock{iterator: iterator()}
	ts := httptest.NewServer(http.HandlerFunc(getInstanceSlices(mock)))
	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte("1")))
	if err != nil {
		t.Fatal(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	var slices []string

	if err := json.Unmarshal(body, &slices); err != nil {
		t.Fatal(err.Error())
	}
	if slices[0] != "2" || slices[1] != "3" {
		t.Error(fmt.Sprintf("unexpected slices returned from server %v", slices))
	}

	defer func() {
		ts.Close()
		res.Body.Close()
	}()
}

func TestIncrClusterId(t *testing.T) {
	mock := clusterMock{iterator: iterator()}
	mock.IncrClusterId()
	mock.IncrClusterId()

	ts := httptest.NewServer(http.HandlerFunc(incrClusterId(mock)))
	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte("1")))
	if err != nil {
		t.Fatal(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(body) != "2" {
		t.Error(fmt.Sprintf("unexpected id returned from server %v", string(body)))
	}

	defer func() {
		ts.Close()
		res.Body.Close()
	}()
}

func startServer() {
	go main()
	runtime.Gosched()
	<-time.After(10 * time.Millisecond) //give the server some time to start
}

func (c clusterMock) GetSlices() []string {
	return []string{"0", "1", "2", "3"}
}

func iterator() func() int {
	i := -1
	return func() int {
		i += 1
		return i
	}
}
func (c clusterMock) IncrClusterId() int {
	return c.iterator()
}
func (c clusterMock) GetInstanceSlices(id int) []string {
	return c.GetSlices()[id*2:(id+1)*2]
}
