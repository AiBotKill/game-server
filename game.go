package main

import (
	"log"
	"time"
)

type game struct {
	Id         string
	State      string
	LastUpdate time.Time
	Bullets    []*bullet
	Players    []*player
}

// NewGame returns a new game with random uuid and with "new" status.
func NewGame() *game {
	g := &game{}
	g.Id = Uuid()
	g.State = "new"
	return g
}

func (g *game) gameOver() bool {
	playerAlive := 0
	for _, p := range g.Players {
		if p.HitPoints > 0 {
			playerAlive++
		}
	}
	return playerAlive < 2
}

// update runs update functions of all entities and bullets.
func (g *game) update(dt time.Duration) {
	g.LastUpdate = time.Now()

	if g.State == "running" {
		for _, b := range g.Bullets {
			b.update(dt)
		}
		for _, p := range g.Players {
			p.update(dt)
		}
		if g.gameOver() {
			g.State = "ended"
		}
	}
}

// TODO: start gameloop
func (g *game) start() {
	// wait for frametime OR until every player has action set
	// update everything
}

// intersectionPoint returns cordinates of intersection between two lines, or nil if they do not collide.
func intersectionPoint(x1, y1, x2, y2, x3, y3, x4, y4 float64) []float64 {
	d := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if d != 0 {
		xi := ((x3-x4)*(x1*y2-y1*x2) - (x1-x2)*(x3*y4-y3*x4)) / d
		yi := ((y3-y4)*(x1*y2-y1*x2) - (y1-y2)*(x3*y4-y3*x4)) / d
		if (x1 < xi && x2 < xi) || (x1 > xi && x2 > xi) ||
			(y1 < yi && y2 < yi) || (y1 > yi && y2 > yi) ||
			(x3 < xi && x4 < xi) || (x3 > xi && x4 > xi) ||
			(y3 < yi && y4 < yi) || (y3 > yi && y4 > yi) {
			return nil
		}
		return []float64{xi, yi}
	}
	return nil
}

type collision struct {
	Entity *player
	Point  [2]float64
}

func (g *game) collision(vector [4]float64) []*collision {
	// Loop trough all entities, optimize opportunity with culling.
	var collisions []*collision
	for _, e := range g.Players {
		// Crate collision borders over entity being checed.
		xmin := e.Location[0] - (e.Dimensions[0] / 2.0)
		ymin := e.Location[1] - (e.Dimensions[1] / 2.0)
		xmax := e.Location[0] + (e.Dimensions[0] / 2.0)
		ymax := e.Location[1] + (e.Dimensions[1] / 2.0)
		borders := [][4]float64{[4]float64{xmin, ymin, xmax, ymin},
			[4]float64{xmin, ymin, xmin, ymax},
			[4]float64{xmax, ymin, xmax, ymax},
			[4]float64{xmin, ymax, xmax, ymax}}

		// Test intersections against all borders
		for _, b := range borders {
			if i := intersectionPoint(vector[0], vector[1], vector[2], vector[3], b[0], b[1], b[2], b[3]); i != nil {
				log.Println(b)
				collisions = append(collisions, &collision{e, [2]float64{i[0], i[1]}})
			}
		}
	}
	return collisions
}

// newBullet adds a new bullet to the bullet updatelist with given location, velocity and damage
func (g *game) newBullet(location [2]float64, velocity Vector, damage float64, shooter *player) {
	log.Println(shooter.entity.Game)
	b := NewBullet(location, velocity, damage, g, shooter)
	g.Bullets = append(g.Bullets, b)
}

func (g *game) newPlayer(location [2]float64, name string) *player {
	p := newPlayer(location, g)
	g.Players = append(g.Players, p)
	return p
}
