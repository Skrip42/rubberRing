package rubberring

type RubberRingStat struct {
	Size             int
	Capacity         int
	ActiveChanks     int
	ActiveCapacity   int
	PassiveChanks    int
	PassiveCapacity  int
	ActiveChanksSize []int
	EndChankNo       int
	StartPosition    int
	EndPosition      int
}

func Stat[V any](ring *RubberRing[V]) RubberRingStat {
	stat := RubberRingStat{
		Size:             ring.size,
		Capacity:         ring.capacity,
		ActiveChanksSize: make([]int, 0, 8),
		PassiveChanks:    len(ring.freeChanks),
		StartPosition:    ring.startPosition,
	}
	st := ring.startChank
	i := 0
	for st != nil {
		if st == ring.endChank {
			stat.EndChankNo = i
		}
		stat.ActiveChanksSize = append(stat.ActiveChanksSize, len(st.data))
		stat.ActiveCapacity += len(st.data)
		stat.ActiveChanks++

		st = st.nextChank
		i++
	}
	stat.PassiveCapacity = ring.capacity - stat.ActiveCapacity
	for i := 0; i < stat.EndChankNo; i++ {
		stat.EndPosition += stat.ActiveChanksSize[i]
	}
	stat.EndPosition += ring.endPosition

	return stat
}
