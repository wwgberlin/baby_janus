package main

import (
	"fmt"
	"net/http"
	"sync"
)

type (
	mutexLocker struct {
		mutex sync.Mutex
	}

	cluster struct {
		mutexLocker
		numInstances int
	}
)

const NUM_PARTS = 136

func (l *mutexLocker) lock(f func()) {
	l.mutex.Lock()
	f()
	l.mutex.Unlock()
}

func newCluster() *cluster {
	return &cluster{mutexLocker: mutexLocker{mutex: sync.Mutex{}}, numInstances: -1}
}

func (c *cluster) incrClusterId(w http.ResponseWriter, r *http.Request) {
	c.lock(func() {
		c.numInstances++
		fmt.Fprintf(w, fmt.Sprintf("%v", c.numInstances))
	})
}
