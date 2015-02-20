package main

import (
	"math"
	"time"
)

type bullet struct {
	Id string `json:"id"`
	Circle
	FiredBy string  `json:"firedBy"`
	Damage  float64 `json:"damage"`
	Dead    bool    `json:"-"`
}

func (b *bullet) update(g *game, dt time.Duration) {
	d := b.Position.Add(b.Velocity.Mul(dt.Seconds()))
	line := Line{A: b.Position, B: d}

	var collisions []*collision

	// Test every tile for collision
	for _, ct := range g.Tiles {
		pCol := ct.Intersect(line)
		for _, c := range pCol {
			col := &collision{}
			col.Collider = b.Id
			col.Target = ct.Id
			col.Position = c
			if !math.IsNaN(col.Position.X) && !math.IsNaN(col.Position.Y) {
				collisions = append(collisions, col)
			}
		}
	}

	// Test every player for collision
	for _, cp := range g.Players {
		if cp.Hitpoints == 0 {
			continue
		}
		if cp.Id != b.FiredBy {
			pCol := cp.Intersect(line)
			for _, c := range pCol {
				col := &collision{}
				col.Collider = b.Id
				col.Target = cp.Id
				col.Position = c
				if !math.IsNaN(col.Position.X) && !math.IsNaN(col.Position.Y) {
					collisions = append(collisions, col)
				}
			}
		}
	}

	// Sort collisions
	for i := 0; i < len(collisions)-1; i++ {
		for j := i + 1; j < len(collisions); j++ {
			v1 := collisions[i].Position.Sub(b.Position)
			v2 := collisions[j].Position.Sub(b.Position)
			if v1.Length() > v2.Length() {
				collisions[i], collisions[j] = collisions[j], collisions[i]
			}
		}
	}

	// If there are collisions, add the closest one to game.collisions
	// and mark bullet as dead.
	if len(collisions) > 0 {
		g.Collisions = append(g.Collisions, collisions[0])
		b.Dead = true
		// Cause damage for players, and add damagemade for the shooter.
		for _, p := range g.Players {
			if p.Id == collisions[0].Target && p.Hitpoints > 0 {
				if p.Hitpoints < b.Damage {
					b.Damage = p.Hitpoints
				}
				p.Hitpoints -= b.Damage
				for _, shooter := range g.Players {
					if shooter.Id == b.FiredBy {
						shooter.DamageMade += b.Damage
						shooter.Killed = append(shooter.Killed, p.Id)
					}
				}
			}
		}
	}
	if math.IsNaN(b.Velocity.Length()) {
		b.Dead = true
	}
	b.Position.X += b.Velocity.X * dt.Seconds()
	b.Position.Y += b.Velocity.Y * dt.Seconds()
}
