package rpc

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestClientLimiter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := NewMockHttpClient(ctrl)
	NewClientLimiter(mockClient)
	size := 11
	requests := make([]*http.Request, 0, size)
	for i := 0; i < size; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:80/req%d", i), nil)
		require.NoError(t, err)
		requests = append(requests, req)

		mockClient.EXPECT().Do(req).Return(&http.Response{
			StatusCode: 200,
		}, nil)
	}
	now := time.Now()
	limiter := NewClientLimiter(mockClient)

	wg := &sync.WaitGroup{}
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter.Do(requests[i])
		}()
	}
	wg.Wait()
	require.True(t, time.Since(now) >= time.Millisecond*time.Duration(millisecondsBetweenRequest*size))
}
