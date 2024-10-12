package rpc

import (
	"net/http"
	"sync"
	"time"
)

const (
	millisecondsBetweenRequest = 110
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ClientLimiter struct {
	m      *sync.Mutex
	time   time.Time
	client HttpClient
}

func NewClientLimiter(client HttpClient) *ClientLimiter {
	return &ClientLimiter{
		m:      &sync.Mutex{},
		time:   time.Now(),
		client: client,
	}
}

func (x *ClientLimiter) Do(req *http.Request) (*http.Response, error) {
	return x.send(req)
}

func (x *ClientLimiter) send(req *http.Request) (*http.Response, error) {
	x.m.Lock()
	defer x.m.Unlock()
	defer func() {
		x.time = time.Now()
	}()
	now := time.Now()
	diff := now.Sub(x.time)
	if diff < time.Millisecond*millisecondsBetweenRequest {
		time.Sleep((time.Millisecond * millisecondsBetweenRequest) - diff)
	}
	res, err := x.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
