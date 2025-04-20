package rubberring

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (s *ConfigSuite) TestDefaultConfig() {
	ring := NewRubberRing[int]()
	stat := ring.Stat()

	s.Equal(256*4, stat.Capacity)
	s.Equal(4, stat.ActiveChanks)
	s.Equal(256, stat.ActiveChanksSize[0])
}

func (s *ConfigSuite) TestWithStartChankCount() {
	tests := []struct {
		name           string
		count          int
		expectedChunks int
	}{
		{
			name:           "positive count",
			count:          5,
			expectedChunks: 5,
		},
		{
			name:           "zero count",
			count:          0,
			expectedChunks: 1, // minimum is 1
		},
		{
			name:           "negative count",
			count:          -1,
			expectedChunks: 1, // minimum is 1
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ring := NewRubberRing[int](WithStartChankCount(tt.count))
			stat := ring.Stat()
			s.Equal(tt.expectedChunks, stat.ActiveChanks)
		})
	}
}

func (s *ConfigSuite) TestWithStartChankSize() {
	tests := []struct {
		name         string
		size         int
		expectedSize int
		expectedCap  int
	}{
		{
			name:         "positive size",
			size:         100,
			expectedSize: 100,
			expectedCap:  100 * 4,
		},
		{
			name:         "zero size",
			size:         0,
			expectedSize: 1, // default size
			expectedCap:  1 * 4,
		},
		{
			name:         "negative size",
			size:         -1,
			expectedSize: 1, // default size
			expectedCap:  1 * 4,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ring := NewRubberRing[int](WithStartChankSize(tt.size))
			stat := ring.Stat()
			s.Equal(tt.expectedSize, stat.ActiveChanksSize[0])
			s.Equal(tt.expectedCap, stat.Capacity)
		})
	}
}

func (s *ConfigSuite) TestWithPassiveChankBufferSize() {
	tests := []struct {
		name                   string
		bufferSize             int
		expectedMaxBufferCount int
		expectedBuffer         int
	}{
		{
			name:                   "positive buffer size",
			bufferSize:             5,
			expectedMaxBufferCount: 5,
			expectedBuffer:         0,
		},
		{
			name:                   "zero buffer size",
			bufferSize:             0,
			expectedMaxBufferCount: 1,
			expectedBuffer:         0,
		},
		{
			name:                   "negative buffer size",
			bufferSize:             -1,
			expectedMaxBufferCount: 1,
			expectedBuffer:         0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ring := NewRubberRing[int](WithPassiveChankBufferSize(tt.bufferSize))
			stat := ring.Stat()
			s.Equal(tt.expectedMaxBufferCount, cap(ring.freeChanks))
			s.Equal(tt.expectedBuffer, stat.PassiveChanks)
		})
	}
}

func (s *ConfigSuite) TestWithGrowStrategy() {
	customStrategy := func(capacity int) (int, int) {
		return 100, 2
	}

	ring := NewRubberRing[int](
		WithGrowStrategy(customStrategy),
		WithStartChankSize(50),
		WithStartChankCount(1),
	)

	// Push more than initial capacity to trigger growth
	for i := 0; i < 150; i++ {
		ring.Push(i)
	}

	stat := ring.Stat()
	s.Equal(250, stat.Capacity) // 100 size * 2 chunks after growth
	s.Equal(3, stat.ActiveChanks)
}

func (s *ConfigSuite) TestMultipleConfigOptions() {
	ring := NewRubberRing[int](
		WithStartChankSize(100),
		WithStartChankCount(2),
		WithPassiveChankBufferSize(5),
	)

	stat := ring.Stat()
	s.Equal(200, stat.Capacity) // 100 * 2
	s.Equal(2, stat.ActiveChanks)
	s.Equal(100, stat.ActiveChanksSize[0])
	s.Equal(0, stat.PassiveChanks)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
