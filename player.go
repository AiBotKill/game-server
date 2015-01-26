package main

import (
	"errors"
	"log"
	"time"
)

type player struct {
	entity
	Name        string
	Speed       float64
	Damage      float64
	BulletSpeed float64
	Looking     Vector
	HitPoints   float64
	LastFired   int
}

func newPlayer(location [2]float64, g *game) *player {
	p := &player{}
	p.Id = Uuid()
	p.HitPoints = 100.0
	p.BulletSpeed = 10.0
	p.Speed = 1.0
	p.Damage = 10.0
	p.Location = location
	p.Velocity = Vector{0, 0}
	p.Dimensions = [2]float64{1, 1}
	p.EntityType = 1
	p.Game = g
	return p
}

// update will update the player
func (p *player) update(dt time.Duration) {
	p.entity.update(dt)
	p.Exhausted = false
}

// action will perform the given action with vector
func (p *player) action(a string, v Vector) error {
	if p.Exhausted {
		return errors.New("player exhausted")
	}
	switch a {
	case "move":
		p.move(v)
	case "look":
		p.look(v)
	case "shoot":
		p.shoot(v)
	}
	p.Exhausted = true
	return nil
}

// shoot will create a new bullet with BulletSpeed velocity, and players location
func (p *player) shoot(at Vector) {
	log.Println(at)
	at.SetLength(p.BulletSpeed)
	log.Println(at)
	p.Game.newBullet(p.Location, at, p.Damage, p)
}

// move sets players velocity torwards the vector, with player.Speed as the maximum velocity
func (p *player) move(at Vector) {
	if at.Length() > p.Speed {
		at.SetLength(p.Speed)
	}
}

// look will set the player.Looking vetor to given vector
func (p *player) look(at Vector) {
	at.Normalize()
	p.Looking = at
}
