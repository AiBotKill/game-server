package main

import "math"

// Circle is a struct of a circle with Radius float64.
type Circle struct {
	Radius   float64 `json:"radius"`
	Position *Vector `json:"position"`
	Velocity *Vector `json:"velocity"`
}

func (circle *Circle) Intersect(l Line) []*Vector {
	E := l.A
	L := l.B
	C := circle.Position
	r := circle.Radius
	d := L.Sub(E)
	f := E.Sub(C)
	a := d.Dot(d)
	b := 2 * f.Dot(d)
	c := f.Dot(f) - r*r

	disc := b*b - 4*a*c

	if disc < 0 {
		// No intersection
		return nil
	}

	disc = math.Sqrt(disc)

	t1 := (-b - disc) / (2 * a)
	t2 := (-b + disc) / (2 * a)

	c1 := &Vector{E.X + d.X*t1, E.Y + d.Y*t1}
	c2 := &Vector{E.X + d.X*t2, E.Y + d.Y*t2}

	t1hit := t1 >= 0 && t1 <= 1
	t2hit := t2 >= 0 && t2 <= 1

	if t1hit && t2hit {
		return []*Vector{c1, c2}
	}

	if t1hit && !t2hit {
		return []*Vector{c1}
	}

	if !t1hit && t2hit {
		return []*Vector{c2}
	}

	return nil
}
