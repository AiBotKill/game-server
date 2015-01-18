package main

import (
	"time"
)

type entity struct {
	Location   [2]float64
	Velocity   [2]float64
	Dimensions [2]float64 // [width, height float64]
}

func (e *entity) update(dt time.Duration) {
	for i := 0; i < len(e.Location); i++ {
		e.Location[i] += e.Velocity[i] * dt.Seconds()
	}
}
