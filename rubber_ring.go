package rubberring

import (
	"io"
	"iter"
)

type GrowStrategy func(capacity int) (newChankSize, newChankCount int)
type SplitStrategy func(capacity, additional_capacity int) int

type chank[V any] struct {
	data      []V
	nextChank *chank[V]
}

type RubberRing[V any] struct {
	startChank    *chank[V]
	startPosition int
	endChank      *chank[V]
	endPosition   int
	freeChanks    chan *chank[V]
	size          int
	capacity      int
	config        config
}

func NewRubberRing[V any](options ...applyConfigFunc) *RubberRing[V] {
	config := defaultConfig
	for _, option := range options {
		option(&config)
	}

	rr := &RubberRing[V]{
		config:     config,
		freeChanks: make(chan *chank[V], config.freeChankBufferSize),
	}
	capacity := config.startChankSize * config.startChankCount
	chanks := createNewChankChain[V](
		config.startChankSize,
		config.startChankCount,
	)
	rr.startChank = chanks
	rr.endChank = chanks
	rr.capacity = capacity

	return rr
}

func (r *RubberRing[V]) Size() int {
	return r.size
}

func (r *RubberRing[V]) Capacity() int {
	return r.capacity
}

func (r *RubberRing[V]) Stat() RubberRingStat {
	return stat(r)
}

func (r *RubberRing[V]) Pull() (V, error) {
	var el V
	if r.size == 0 {
		return el, io.EOF
	}
	el = r.startChank.data[r.startPosition]
	r.startPosition++
	r.size--
	if r.startPosition >= len(r.startChank.data) {
		newStartChank := r.startChank.nextChank
		r.startChank.nextChank = nil
		select {
		case r.freeChanks <- r.startChank:
		default:
			r.capacity -= len(r.startChank.data)
		}
		r.startChank = newStartChank
		r.startPosition = 0
	}
	return el, nil
}

func (r *RubberRing[V]) Push(el V) {
	r.endChank.data[r.endPosition] = el
	r.endPosition++
	r.size++
	if r.endPosition >= len(r.endChank.data) {
		var newEndChank *chank[V]
		if r.endChank.nextChank != nil {
			newEndChank = r.endChank.nextChank
		} else {
			select {
			case newEndChank = <-r.freeChanks:
			default:
				newChankSize, newChankCount := r.config.growStrategy(r.capacity)
				newChanks := createNewChankChain[V](
					newChankSize, newChankCount,
				)
				newEndChank = newChanks
				r.capacity += newChankSize * newChankCount
			}
		}
		r.endChank.nextChank = newEndChank
		r.endChank = newEndChank
		r.endPosition = 0
	}
}

func (r *RubberRing[V]) Elements() iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			v, err := r.Pull()
			if err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func createNewChankChain[V any](
	chankSize int,
	chankCount int,
) *chank[V] {
	var chk *chank[V]
	for i := 0; i < chankCount; i++ {
		newChank := &chank[V]{
			data:      make([]V, chankSize),
			nextChank: chk,
		}
		chk = newChank
	}
	return chk
}
