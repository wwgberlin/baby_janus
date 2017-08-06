package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestBasicProxy(t *testing.T) {
	calledBack := false
	origin := "/origin"
	target := "http://127.0.0.1:8080/target"

	http.HandleFunc("/origin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "it worked")
		calledBack = true
	})

	go main()

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
