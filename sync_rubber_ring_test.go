package rubberring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type SyncRubberRingSuite struct {
	suite.Suite
	ring *SyncRubberRing[int]
}

func TestSyncRubberRingSuite(t *testing.T) {
	suite.Run(t, &SyncRubberRingSuite{})
}

func (s *SyncRubberRingSuite) SetupTest() {
	s.ring = NewSyncRubberRing[int](
		WithStartChankSize(3),
		WithStartChankCount(2),
	)
}

func (s *SyncRubberRingSuite) TearDownTest() {
	s.NoError(goleak.Find())
	s.ring = nil
}

func (s *SyncRubberRingSuite) TestBasicOperations() {
	ctx := context.Background()

	// Test push and pull
	s.ring.Push(1)
	s.ring.Push(2)
	s.ring.Push(3)

	s.Equal(3, s.ring.Size())
	s.Equal(6, s.ring.Capacity())

	val, err := s.ring.Pull(ctx)
	s.NoError(err)
	s.Equal(1, val)

	val, err = s.ring.Pull(ctx)
	s.NoError(err)
	s.Equal(2, val)

	val, err = s.ring.Pull(ctx)
	s.NoError(err)
	s.Equal(3, val)

	s.Equal(0, s.ring.Size())
	s.Equal(6, s.ring.Capacity())
}

func (s *SyncRubberRingSuite) TestContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to pull
	pullDone := make(chan struct{})
	go func() {
		_, err := s.ring.Pull(ctx)
		s.Equal(context.Canceled, err)
		close(pullDone)
	}()

	// Give some time for the goroutine to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context
	cancel()

	// Wait for the goroutine to finish
	<-pullDone
}

func (s *SyncRubberRingSuite) TestConcurrentAccess() {
	ctx := context.Background()
	const numGoroutines = 5
	const itemsPerGoroutine = 100

	// Start producer goroutines
	for i := range numGoroutines {
		go func(start int) {
			for j := range itemsPerGoroutine {
				s.ring.Push(start + j)
			}
		}(i * itemsPerGoroutine)
	}

	// Start consumer goroutines
	consumed := make(chan int, numGoroutines*itemsPerGoroutine)
	for range numGoroutines {
		go func() {
			for range itemsPerGoroutine {
				val, err := s.ring.Pull(ctx)
				s.NoError(err)
				consumed <- val
			}
		}()
	}

	// Collect all consumed values
	values := make([]int, 0, numGoroutines*itemsPerGoroutine)
	for range numGoroutines * itemsPerGoroutine {
		values = append(values, <-consumed)
	}

	// Verify all values were consumed
	s.Equal(numGoroutines*itemsPerGoroutine, len(values))
	s.Equal(0, s.ring.Size())
}

func (s *SyncRubberRingSuite) TestStat() {
	s.ring.Push(1)
	s.ring.Push(2)
	s.ring.Push(3)

	stat := s.ring.Stat()

	s.Equal(3, stat.Size)
	s.Equal(6, stat.Capacity)
	s.Equal(2, stat.ActiveChanks)
	s.Equal(6, stat.ActiveCapacity)
	s.Equal(0, stat.PassiveChanks)
	s.Equal(0, stat.PassiveCapacity)
	s.Len(stat.ActiveChanksSize, 2)
	s.Equal([]int{3, 3}, stat.ActiveChanksSize)
}

func (s *SyncRubberRingSuite) TestElements() {
	ctx, canceled := context.WithCancel(context.Background())

	var result []int
	go func() {
		// Test iterator
		for v := range s.ring.Elements(ctx) {
			result = append(result, v)
		}
	}()

	// Push some values
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.ring.Push(v)
	}
	time.Sleep(time.Millisecond * 100)
	canceled()

	s.Equal(values, result)
	s.Equal(0, s.ring.Size())
}

func (s *SyncRubberRingSuite) TestElementsWithContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())

	// Push some values
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.ring.Push(v)
	}

	// Start iterator in a goroutine
	iterDone := make(chan struct{})
	go func() {
		var result []int
		for v := range s.ring.Elements(ctx) {
			result = append(result, v)
			if len(result) == 2 {
				cancel()
				break
			}
		}
		close(iterDone)
	}()

	// Wait for the iterator to finish
	<-iterDone

	// Verify that not all values were consumed
	s.Greater(s.ring.Size(), 0)
}
