package rubberring

type config struct {
	startChankSize        int
	pasiveChankBufferSize int
	startChankCount       int
	growStrategy          GrowStrategy
}

var defaultConfig = config{
	startChankSize:        256,
	startChankCount:       4,
	pasiveChankBufferSize: 3,
	growStrategy:          func(capacity int) (int, int) { return 256, 4 },
}

type applyConfigFunc func(o *config)

func WithStartChankCount(count int) applyConfigFunc {
	if count < 1 {
		count = 1
	}
	return func(c *config) {
		c.startChankCount = count
	}
}

func WithStartChankSize(size int) applyConfigFunc {
	if size < 1 {
		size = 1
	}
	return func(c *config) {
		c.startChankSize = size
	}
}

func WithPassiveChankBufferSize(bufferSize int) applyConfigFunc {
	if bufferSize < 1 {
		bufferSize = 1
	}
	return func(c *config) {
		c.pasiveChankBufferSize = bufferSize
	}
}

func WithGrowStrategy(strategy GrowStrategy) applyConfigFunc {
	return func(c *config) {
		c.growStrategy = strategy
	}
}
