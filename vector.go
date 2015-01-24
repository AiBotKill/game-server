package main

import (
	"math"
)

// Vector is a [2]float64 type with extra methods for calculating and
// manipulating length and angle.
type Vector [2]float64

// Normalize sets the length of vector to 1.
func (v *Vector) Normalize() {
	l := v.Length()
	if l != 0 {
		v[0] = v[0] / l
		v[1] = v[1] / l
	}
}

// Length returns the vector's length.
func (v *Vector) Length() float64 {
	return math.Sqrt(math.Pow(v[0], 2) + math.Pow(v[1], 2))
}

// Angle returns angle in radians.
func (v *Vector) Angle() float64 {
	return math.Atan2(v[1], v[0])
}

// SetLength sets the length of vector so that the angle stays the same.
func (v *Vector) SetLength(l float64) {
	v.Normalize()
	v[0] *= l
	v[1] *= l
}

// SetAngle sets the angle of vector so that the length stays the same.
func (v *Vector) SetAngle(angle float64) {
	l := v.Length()
	v[0] = l * math.Cos(angle)
	v[1] = l * math.Sin(angle)
}

// Rotate rotates the vector, keeping length the same.
func (v *Vector) Rotate(angle float64) {
	v.SetAngle(v.Angle() + angle)
}
