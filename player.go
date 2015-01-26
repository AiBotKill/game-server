package main

type player struct {
	entity
	Name        string
	Speed       float64
	Damage      float64
	BulletSpeed float64
}

func newPlayer() *player {
	p := &player{}
	p.entity = entity{}
	p.HitPoints = 100.0
	p.Speed = 1.0
	p.Damage = 10.0
	return p
}

func (p *player) shoot(at Vector) {
	at.SetLength(p.BulletSpeed)
	p.Game.newBullet(p.Location, at, p.Damage, p)
}

func (p *player) move(at Vector) {
	if at.Length() > p.Speed {
		at.SetLength(p.Speed)
	}
	p.Velocity = at
}

func (p *player) look(at Vector) {
	at.Normalize()
	p.Looking = at
}
