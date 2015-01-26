package main

import (
	"time"
)

type Entity interface {
	GetLocation() [2]float64
	GetDimensions() [2]float64
	GetVelocity() Vector
	StopsBullets() bool
	StopsVision() bool
	Collision(e *Entity)
}

func NewEntity(location [2]float64, entityType int, g *game) *entity {
	e := &entity{}
	e.Id = Uuid()
	e.Location = location
	e.Velocity = Vector{0, 0}
	e.Dimensions = [2]float64{1, 1}
	e.EntityType = entityType
	e.Game = g
	return e
}

type entity struct {
	Id         string
	Game       *game
	Location   [2]float64
	Velocity   Vector
	Dimensions [2]float64 // [width, height float64]
	EntityType int
	Exhausted  bool
}

// GetLocation to implement Entity interface
func (e *entity) GetLocation() [2]float64 {
	return e.Location
}

// GetDimensions to implement Entity interface
func (e *entity) GetDimensions() [2]float64 {
	return e.Dimensions
}

// GetVelocity to implement Entity interface
func (e *entity) GetVelocity() Vector {
	return e.Velocity
}

// Update will update location according to velocity
func (e *entity) update(dt time.Duration) {
	for i := 0; i < len(e.Location); i++ {
		e.Location[i] += e.Velocity[i] * dt.Seconds()
	}
	e.Exhausted = false
}
