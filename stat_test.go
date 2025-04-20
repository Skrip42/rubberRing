package rubberring

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatSuite struct {
	suite.Suite
	ring *RubberRing[int]
}

func (s *StatSuite) SetupTest() {
	s.ring = NewRubberRing[int](
		WithStartChankSize(3),
		WithStartChankCount(2),
	)
}

func (s *StatSuite) TearDownTest() {
	s.ring = nil
}

func (s *StatSuite) TestInitialState() {
	stat := s.ring.Stat()

	s.Equal(0, stat.Size)
	s.Equal(6, stat.Capacity) // 3 * 2
	s.Equal(2, stat.ActiveChanks)
	s.Equal(6, stat.ActiveCapacity)
	s.Equal(0, stat.PassiveChanks)
	s.Equal(0, stat.PassiveCapacity)
	s.Len(stat.ActiveChanksSize, 2)
	s.Equal([]int{3, 3}, stat.ActiveChanksSize)
	s.Equal(0, stat.EndChankNo)
	s.Equal(0, stat.StartPosition)
	s.Equal(0, stat.EndPosition)
}

func (s *StatSuite) TestStatAfterPush() {
	// Push some elements
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
	s.Equal(1, stat.EndChankNo)
	s.Equal(0, stat.StartPosition)
	s.Equal(3, stat.EndPosition)
}

func (s *StatSuite) TestStatAfterPull() {
	// Push and then pull some elements
	s.ring.Push(1)
	s.ring.Push(2)
	s.ring.Push(3)
	s.ring.Pull()
	s.ring.Pull()

	stat := s.ring.Stat()

	s.Equal(1, stat.Size)
	s.Equal(6, stat.Capacity)
	s.Equal(2, stat.ActiveChanks)
	s.Equal(6, stat.ActiveCapacity)
	s.Equal(0, stat.PassiveChanks)
	s.Equal(0, stat.PassiveCapacity)
	s.Len(stat.ActiveChanksSize, 2)
	s.Equal([]int{3, 3}, stat.ActiveChanksSize)
	s.Equal(1, stat.EndChankNo)
	s.Equal(2, stat.StartPosition)
	s.Equal(3, stat.EndPosition)
}

func (s *StatSuite) TestStatWithGrowth() {
	// Push more elements than initial capacity to trigger growth
	for i := 0; i < 8; i++ {
		s.ring.Push(i)
	}

	stat := s.ring.Stat()

	s.Equal(8, stat.Size)
	s.Greater(stat.Capacity, 6)       // Should be greater than initial capacity
	s.Greater(stat.ActiveChanks, 2)   // Should have more than initial chunks
	s.Greater(stat.ActiveCapacity, 6) // Should have more than initial capacity
	s.Equal(0, stat.PassiveChanks)
	s.Equal(0, stat.PassiveCapacity)
	s.Greater(len(stat.ActiveChanksSize), 2)
	s.Equal(2, stat.EndChankNo)
	s.Equal(0, stat.StartPosition)
	s.Equal(8, stat.EndPosition)
}

func (s *StatSuite) TestStatEmptyRing() {
	// Pull all elements
	s.ring.Push(1)
	s.ring.Push(2)
	s.ring.Pull()
	s.ring.Pull()

	stat := s.ring.Stat()

	s.Equal(0, stat.Size)
	s.Equal(6, stat.Capacity)
	s.Equal(2, stat.ActiveChanks)
	s.Equal(6, stat.ActiveCapacity)
	s.Equal(0, stat.PassiveChanks)
	s.Equal(0, stat.PassiveCapacity)
	s.Len(stat.ActiveChanksSize, 2)
	s.Equal([]int{3, 3}, stat.ActiveChanksSize)
	s.Equal(0, stat.EndChankNo)
	s.Equal(2, stat.StartPosition)
	s.Equal(2, stat.EndPosition)
}

func TestStatSuite(t *testing.T) {
	suite.Run(t, new(StatSuite))
}
