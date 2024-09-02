package tests

import (
	"google.golang.org/grpc/credentials/insecure"
	"time"

	"google.golang.org/grpc"
)

func WaitForPort(address string) bool {
	waitChan := make(chan struct{})

	go func() {
		for {
			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if conn != nil {
				waitChan <- struct{}{}
				return
			}
		}
	}()

	timeout := time.After(3 * time.Second)
	select {
	case <-waitChan:
		return true
	case <-timeout:
		return false
	}
}
