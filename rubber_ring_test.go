package rubberring

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RubberRingSuite struct {
	suite.Suite
	ring *RubberRing[int]
}

func (s *RubberRingSuite) SetupTest() {
	s.ring = NewRubberRing[int]()
}

func (s *RubberRingSuite) TearDownTest() {
	s.ring = nil
}

func (s *RubberRingSuite) TestNewRubberRing() {
	tests := []struct {
		name           string
		options        []applyConfigFunc
		expectedSize   int
		expectedCap    int
		expectedChunks int
	}{
		{
			name:           "default configuration",
			options:        nil,
			expectedSize:   0,
			expectedCap:    256 * 4,
			expectedChunks: 4,
		},
		{
			name: "custom configuration",
			options: []applyConfigFunc{
				WithStartChankSize(100),
				WithStartChankCount(2),
			},
			expectedSize:   0,
			expectedCap:    100 * 2,
			expectedChunks: 2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rr := NewRubberRing[int](tt.options...)
			s.Equal(tt.expectedSize, rr.Size())
			s.Equal(tt.expectedCap, rr.Capacity())
		})
	}
}

func (s *RubberRingSuite) TestPushPull() {
	// Test empty ring
	_, err := s.ring.Pull()
	s.Equal(io.EOF, err)

	// Test basic push and pull
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.ring.Push(v)
	}

	s.Equal(len(values), s.ring.Size())

	for _, want := range values {
		got, err := s.ring.Pull()
		s.NoError(err)
		s.Equal(want, got)
	}

	s.Equal(0, s.ring.Size())
}

func (s *RubberRingSuite) TestCapacityGrowth() {
	rr := NewRubberRing[int](
		WithStartChankSize(2),
		WithStartChankCount(1),
	)

	initialCap := rr.Capacity()

	// Push more items than initial capacity
	for i := 0; i < 10; i++ {
		rr.Push(i)
	}

	s.Greater(rr.Capacity(), initialCap)
}

func (s *RubberRingSuite) TestElements() {
	values := []int{1, 2, 3, 4, 5}

	for _, v := range values {
		s.ring.Push(v)
	}

	var result []int
	for v := range s.ring.Elements() {
		result = append(result, v)
	}

	s.Equal(len(values), len(result))
	s.Equal(values, result)
}

func (s *RubberRingSuite) TestStat() {
	rr := NewRubberRing[int](
		WithStartChankSize(2),
		WithStartChankCount(1),
	)

	stat := rr.Stat()

	s.Equal(0, stat.Size)
	s.Equal(2, stat.Capacity)
	s.Equal(1, stat.ActiveChanks)
	s.Len(stat.ActiveChanksSize, 1)
}

func TestRubberRingSuite(t *testing.T) {
	suite.Run(t, new(RubberRingSuite))
}
