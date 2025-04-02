package rubberring

import (
	"io"
)

type GrowStrategy func(size int) int

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
	chanks, capacity := createNewChankChain[V](config.startCapacity, config.splitFactor)
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
				newChanks, newChanksCapacity := createNewChankChain[V](
					r.config.growStrategy(r.capacity),
					r.config.splitFactor,
				)
				newEndChank = newChanks
				r.capacity += newChanksCapacity
			}
		}
		r.endChank.nextChank = newEndChank
		r.endChank = newEndChank
		r.endPosition = 0
	}
}

func createNewChankChain[V any](
	size int,
	splitFactor int,
) (*chank[V], int) {
	var chk *chank[V]
	chankSize := size / splitFactor
	for i := 0; i < splitFactor; i++ {
		newChank := &chank[V]{
			data:      make([]V, chankSize),
			nextChank: chk,
		}
		chk = newChank
	}
	return chk, chankSize * splitFactor
}
