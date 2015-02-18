package main

import (
	"log"
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
	// Parse the action
	switch p.Action.Type {
	case "move":
		p.Velocity = p.Action.Direction
		if p.Velocity.Length() > MAX_SPEED {
			p.Velocity = p.Velocity.Normalize().Mul(MAX_SPEED)
		}
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

	// Test every tile for collision
	var tCols []*collision
	for _, ct := range g.Tiles {
		tCol := ct.Intersect(line)
		for _, c := range tCol {
			col := &collision{
				Collider: p.Id,
				Target:   ct.Id,
				Position: c,
			}
			tCols = append(tCols, col)
		}
	}
	tCols = SortCollisions(tCols, p.Position)
	if len(tCols) > 1 {
		// TODO : This is very sticky collision, might be difficult to handle for AI.
		p.Position = tCols[0].Position
	} else {
		p.Position.X += p.Velocity.X * dt.Seconds()
		p.Position.Y += p.Velocity.Y * dt.Seconds()
	}

}
