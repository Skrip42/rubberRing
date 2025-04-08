package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	rubberring "github.com/Skrip42/rubberRing"
)

func main() {
	// This version of the buffer is thread-safe.
	// Also allows you to wait for an element to appear if the buffer is empty.
	// To terminate the wait urgently, you need to close the context.
	rr := rubberring.NewSyncRubberRing[int](
		rubberring.WithStartChankCount(2),
		rubberring.WithStartChankSize(2),
		rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }),
		rubberring.WithFreeChankBufferSize(2),
	)
	resultFromRootine1 := []int{}
	resultFromRootine2 := []int{}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		for i := 0; i < 12; i++ {
			rr.Push(i)
			time.Sleep(100 * time.Millisecond)

		}
		wg.Done()
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	go func() {
		defer wg.Done()
		for ctx.Err() == nil {
			v, err := rr.Pull(ctx)
			if err != nil {
				return
			}
			resultFromRootine1 = append(resultFromRootine1, v)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for ctx.Err() == nil {
			v, err := rr.Pull(ctx)
			if err != nil {
				return
			}
			resultFromRootine2 = append(resultFromRootine2, v)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	wg.Wait()

	// The exact composition is not predetermined,
	// but it should turn out something similar to
	// [1 3 5 7 9 11]
	// [0 2 4 6 8 10]
	fmt.Println(resultFromRootine1)
	fmt.Println(resultFromRootine2)
}
