package main

import (
	"context"
	"fmt"
	"time"

	rubberring "github.com/Skrip42/rubberRing"
)

func main() {
	// Create simple rubberRing buffer
	rr := rubberring.NewRubberRing[int](
		rubberring.WithStartChankCount(2),
		rubberring.WithStartChankSize(2),
		rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }),
		rubberring.WithFreeChankBufferSize(2),
	)

	// push some data
	for i := 0; i < 8; i++ {
		rr.Push(i)
	}

	// get all elements from boffer
	for v := range rr.Elements() {
		fmt.Println(v)
	}

	// Create sync rubberRing buffer
	srr := rubberring.NewSyncRubberRing[int](
		rubberring.WithStartChankCount(2),
		rubberring.WithStartChankSize(2),
		rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }),
		rubberring.WithFreeChankBufferSize(2),
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// push some data
		for i := 0; i < 8; i++ {
			srr.Push(i)
			time.Sleep(100 * time.Millisecond)
		}
		cancel()
	}()

	// get all the elements from buffer as they come in
	for v := range srr.Elements(ctx) {
		fmt.Println(v)
	}
}
