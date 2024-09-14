package rpc

import (
	"net/http"
	"sync"
	"time"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ClientLimiter struct {
	m      *sync.Mutex
	time   time.Time
	client *http.Client
}

func NewClientLimiter(client *http.Client) *ClientLimiter {
	return &ClientLimiter{
		m:      &sync.Mutex{},
		time:   time.Now(),
		client: client,
	}
}

type AsyncRes struct {
	HttpResponse *http.Response
	Err          error
}

func (x *ClientLimiter) Do(req *http.Request) (*http.Response, error) {
	stream := x.send(req)
	asyncRes := <-stream
	return asyncRes.HttpResponse, asyncRes.Err
}

func (x *ClientLimiter) send(req *http.Request) <-chan AsyncRes {
	x.m.Lock()
	defer x.m.Unlock()
	now := time.Now()
	diff := now.Sub(x.time)
	if diff < time.Millisecond*110 {
		time.Sleep((time.Millisecond * 100) - diff)
	}
	var httpResStream chan AsyncRes
	go func() {
		defer close(httpResStream)
		res, err := x.client.Do(req)
		httpResStream <- AsyncRes{
			HttpResponse: res,
			Err:          err,
		}
	}()
	x.time = time.Now()
	return httpResStream
}
