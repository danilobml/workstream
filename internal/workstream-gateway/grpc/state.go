package grpc

import (
	"errors"
	"sync"

	"google.golang.org/grpc"
)

var (
	mu   sync.RWMutex
	conn *grpc.ClientConn
	dial error
)

func SetClient(c *grpc.ClientConn, err error) {
	mu.Lock()
	defer mu.Unlock()
	conn = c
	dial = err
}

func GetClient() (*grpc.ClientConn, error) {
	mu.RLock()
	defer mu.RUnlock()

	if dial != nil {
		return nil, dial
	}
	if conn == nil {
		return nil, errors.New("grpc client connection is nil")
	}
	return conn, nil
}
