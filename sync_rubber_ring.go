package rubberring

import (
	"context"
	"iter"
	"sync"

	syncutils "github.com/Skrip42/syncUtils"
)

type SyncRubberRing[V any] struct {
	ring *RubberRing[V]
	cond *syncutils.Cond
	mu   *sync.Mutex
}

func NewSyncRubberRing[V any](options ...applyConfigFunc) *SyncRubberRing[V] {
	return &SyncRubberRing[V]{
		ring: NewRubberRing[V](options...),
		cond: syncutils.NewCond(),
		mu:   &sync.Mutex{},
	}
}

func (r *SyncRubberRing[V]) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Size()
}

func (r *SyncRubberRing[V]) Capacity() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Capacity()
}

func (r *SyncRubberRing[V]) Stat() RubberRingStat {
	r.mu.Lock()
	defer r.mu.Unlock()
	return stat(r.ring)
}

func (r *SyncRubberRing[V]) Push(value V) {
	r.mu.Lock()
	r.ring.Push(value)
	r.cond.Signal()
	r.mu.Unlock()
}

func (r *SyncRubberRing[V]) Pull(ctx context.Context) (V, error) {
	var v V
	r.mu.Lock()
	for {
		if r.ring.Size() > 0 {
			v, err := r.ring.Pull()
			r.mu.Unlock()
			return v, err
		}
		wait := r.cond.Wait()
		r.mu.Unlock()
		select {
		case <-wait:
		case <-ctx.Done():
			return v, ctx.Err()
		}
		r.mu.Lock()
	}
}

func (r *SyncRubberRing[V]) Elements(ctx context.Context) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			v, err := r.Pull(ctx)
			if err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}
