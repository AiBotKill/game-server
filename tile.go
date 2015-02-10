package main

import "time"

type tile struct {
	Id string `json:"id"`
	Rectangle
}

func (t *tile) update(g *game, dt time.Duration) {
	t.Position.X += t.Velocity.X * dt.Seconds()
	t.Position.Y += t.Velocity.Y * dt.Seconds()
}
