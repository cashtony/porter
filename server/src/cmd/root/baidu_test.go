package main

import (
	"sync"
	"testing"
)

func TestNewBaiduUser(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 1; i < 100; i++ {
		wg.Add(1)
		go func() {
			_, err := NewBaiduUser("0yUjhJQnlEQjZHRmJOQ2dtbmtoRn5xWHo4a3JlMEtieFhRdndIOWV3MVV2ZWxmRVFBQUFBJCQAAAAAAAAAAAEAAABDV7QWeG4xMjEzMDAxOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFQwwl9UMMJfRE")
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()

	}
	wg.Wait()
}
