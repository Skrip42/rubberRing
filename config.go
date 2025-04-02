package rubberring

type config struct {
	startCapacity       int
	freeChankBufferSize int
	growStrategy        GrowStrategy
	splitFactor         int
}

var defaultConfig = config{
	startCapacity:       24,
	freeChankBufferSize: 2,
	growStrategy:        func(size int) int { return 24 },
	splitFactor:         3,
}

type applyConfigFunc func(o *config)

func WithStartCapacity(capacity int) applyConfigFunc {
	return func(c *config) {
		c.startCapacity = capacity
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

func WithSplitFactor(factor int) applyConfigFunc {
	return func(c *config) {
		c.splitFactor = factor
	}
}
