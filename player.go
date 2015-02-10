package main

import "time"

type player struct {
	Id string `json:"id"`
	Circle
	Name       string  `json:"name"`
	LookingAt  *Vector `json:"lookingAt"`
	Hitpoints  float64 `json:"hitpoints"`
	DamageMade float64 `json:"damageMade"`
	Action     struct {
		Type      string  `json:"type"`
		Direction *Vector `json:"direction"`
	} `json:"action"`
}

func (p *player) update(g *game, dt time.Duration) {
	d := p.Position.Add(p.Velocity.Mul(dt.Seconds()))
	line := Line{A: p.Position, B: d}
	// Test every player for collision
	for _, cp := range g.Players {
		if cp != p {
			pCol := cp.Intersect(line)
			for _, c := range pCol {
				col := &collision{}
				col.Collider = p.Id
				col.Target = cp.Id
				col.Position = c
				g.Collisions = append(g.Collisions, col)
			}
		}
	}
	p.Position.X += p.Velocity.X * dt.Seconds()
	p.Position.Y += p.Velocity.Y * dt.Seconds()
}
