package rubberring

type config struct {
	startChankSize      int
	freeChankBufferSize int
	startChankCount     int
	growStrategy        GrowStrategy
}

var defaultConfig = config{
	startChankSize:      256,
	startChankCount:     4,
	freeChankBufferSize: 3,
	growStrategy:        func(capacity int) (int, int) { return 256, 4 },
}

type applyConfigFunc func(o *config)

func WithStartChankCount(count int) applyConfigFunc {
	return func(c *config) {
		c.startChankCount = count
	}
}

func WithStartChankSize(size int) applyConfigFunc {
	return func(c *config) {
		c.startChankSize = size
	}
}

func WithFreeChankBufferSize(bufferSize int) applyConfigFunc {
	return func(c *config) {
		c.freeChankBufferSize = bufferSize
	}
}

func WithGrowStrategy(strategy GrowStrategy) applyConfigFunc {
	return func(c *config) {
		c.growStrategy = strategy
	}
}
