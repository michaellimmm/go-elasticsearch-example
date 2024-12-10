package sampler

import (
	"math/rand"
	"time"
)

type ResevoirSampler[T any] struct {
	reservoir   []T
	count       int
	totalSample int
	r           *rand.Rand
}

// NewReservoirSampler creates a new ReservoirSampler[T] with the given sample size k.
// The reservoir is initially empty and the count is 0.
// The rand.Rand is used to generate random numbers.
func NewReservoirSampler[T any](k int) *ResevoirSampler[T] {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &ResevoirSampler[T]{
		reservoir:   make([]T, 0, k),
		count:       0,
		totalSample: k,
		r:           r,
	}
}

// Add adds a new item to the sampler.
//
// If the reservoir is not full (i.e., its size is less than totalSample), the item is appended to the reservoir.
// Otherwise, a random index j is generated between 0 and count-1. If j is less than totalSample, the item at index j
// is replaced with the new item.
func (rs *ResevoirSampler[T]) Add(item T) {
	rs.count++

	if len(rs.reservoir) < rs.totalSample {
		rs.reservoir = append(rs.reservoir, item)
		return
	}

	j := rs.r.Intn(rs.count)
	if j < rs.totalSample {
		rs.reservoir[j] = item
	}
}

func (rs *ResevoirSampler[T]) GetSample() []T {
	return rs.reservoir
}

func (rs *ResevoirSampler[T]) GetTotalSample() int {
	return rs.totalSample
}
