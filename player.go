package main

import (
	"log"
	"math"
	"time"
)

const (
	RATE_OF_FIRE  = time.Millisecond * 500
	MAX_SPEED     = 1.0
	BULLET_SPEED  = 5.0
	BULLET_DAMAGE = 10.0
)

type player struct {
	Circle
	Id         string    `json:"id"`
	BotId      string    `json:"botId"`
	Name       string    `json:"name"`
	Team       int64     `json:"team"`
	LookingAt  *Vector   `json:"lookingAt"`
	Hitpoints  float64   `json:"hitpoints"`
	DamageMade float64   `json:"damageMade"`
	Killed     []string  `json:"killed"`
	LastFired  time.Time `json:"lastFired"`
	Linkdead   bool      `json:"linkdead"`
	Action     struct {
		Type      string  `json:"type"`
		Direction *Vector `json:"direction"`
	} `json:"action"`
}

func (p *player) update(g *game, dt time.Duration) {
	if p.Hitpoints == 0 {
		return
	}
	// Parse the action
	switch p.Action.Type {
	case "move":
		acDir := &Vector{p.Action.Direction.X, p.Action.Direction.Y}
		acDir.Sub(p.Position)
		if acDir.Length() > MAX_SPEED {
			acDir = acDir.Normalize().Mul(MAX_SPEED)
		}
		p.Velocity = acDir
		log.Println(p.Id + " moving")
	case "look":
		p.LookingAt = p.Action.Direction
		log.Println(p.Id + " looking")
	case "shoot":
		p.LookingAt = p.Action.Direction
		if p.LastFired.Add(RATE_OF_FIRE).Before(g.LastUpdate.Add(dt)) {
			bulletPos := &Vector{p.Position.X, p.Position.Y}
			b := g.newBullet(bulletPos, p.Action.Direction.Normalize().Mul(BULLET_SPEED), p.Id)
			b.Damage = 10.0
			p.LastFired = g.LastUpdate.Add(dt)
			log.Println(p.Id + " shooting succesfully")
		}
		log.Println(p.Id + " shooting")
	}

	d := p.Position.Add(p.Velocity.Mul(dt.Seconds()))
	line := Line{A: p.Position, B: d}

	var collisions []*collision

	// Test every tile for collision
	for _, ct := range g.Tiles {
		pCol := ct.Intersect(line)
		for _, c := range pCol {
			col := &collision{}
			col.Collider = p.Id
			col.Target = ct.Id
			col.Position = c
			if !math.IsNaN(col.Position.X) && !math.IsNaN(col.Position.Y) {
				collisions = append(collisions, col)
			}
		}
	}

	// Sort collisions
	for i := 0; i < len(collisions)-1; i++ {
		for j := i + 1; j < len(collisions); j++ {
			v1 := collisions[i].Position.Sub(p.Position)
			v2 := collisions[j].Position.Sub(p.Position)
			if v1.Length() > v2.Length() {
				collisions[i], collisions[j] = collisions[j], collisions[i]
			}
		}
	}

	if len(collisions) > 0 {
		dist := collisions[0].Position.Sub(p.Position).Length()
		pVec := collisions[0].Position.Sub(p.Position)
		p.Position = p.Position.Add(pVec.Normalize().Mul(dist - p.Radius))
	} else {
		if math.IsNaN(p.Velocity.Length()) {
			p.Velocity = &Vector{0, 0}
		}
		p.Position.X += p.Velocity.X * dt.Seconds()
		p.Position.Y += p.Velocity.Y * dt.Seconds()
	}

}
