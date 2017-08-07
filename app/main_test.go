package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"runtime"
)

func startServer() {
	go main()
	runtime.Gosched()
}

func TestBasicProxy(t *testing.T) {
	startServer()

	calledBack := false
	origin := "/origin"
	target := "http://127.0.0.1:8080/target"

	http.HandleFunc("/target", func(w http.ResponseWriter, r *http.Request) {
		calledBack = true
	})

	if b, err := json.Marshal(Endpoint{Origin: origin, Target: target}); err != nil {
		t.Error(err.Error())
		return
	} else if _, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/register_endpoint"), "application/json", bytes.NewBuffer(b)); err != nil {
		t.Error(err.Error())
	}

	if _, err := http.Get("http://127.0.0.1:8080/origin"); err != nil {
		t.Error(err.Error())
	}
	if !calledBack {
		t.Error("didn't work")
	}
}

func TestIncrementClusterId(t *testing.T) {
	go incrementClusterId()
}
