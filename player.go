package main

type player struct {
	entity
	Name   string
	Speed  float64
	Damage float64
}

func newPlayer() *player {
	p := &player{}
	p.entity = entity{}
	p.HitPoints = 100.0
	p.Speed = 1.0
	p.Damage = 10.0
	return p
}

func (p *player) shoot(at [2]float64) {
	var velocity [2]float64
	velocity[0] = 10.0
	velocity[1] = 0.0
	p.Game.newBullet(p.Location, velocity, 10.0, p)
}

func (p *player) move(at [2]float64) {}

func (p *player) look(at [2]float64) {}
