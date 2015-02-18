package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type gameStateMsg struct {
	Type       string    `json:"type"`
	Id         string    `json:"id"`
	StartTime  time.Time `json:"startTime`
	TimeLeftMs float64   `json:"timeLeftMs"`
	State      string    `json:"state"`
	Players    []*player `json:"players"`
	Bullets    []*bullet `json:"bullets"`
	Tiles      []*tile   `json:"tiles"`
}

type collision struct {
	Collider string  `json:"collider"`
	Target   string  `json:"target"`
	Position *Vector `json:"position"`
}

type game struct {
	Id                string        `json:"id"`
	StartTime         time.Time     `json:"startTime"`
	TimeLimit         time.Duration `json:"timeLimit"`
	State             string        `json:"state"`
	GameArea          [2]float64    `json:"gameArea"`
	Mode              string        `json:"mode"`
	Environment       string        `json:"environment"`
	Tiles             []*tile       `json:"tiles"`
	Players           []*player     `json:"players"`
	Bullets           []*bullet     `json:"bullets"`
	Collisions        []*collision  `json:"collisions"`
	StartingPositions []*Vector     `json:"startingPositions"`
	LastUpdate        time.Time     `json:"-"`
}

func newGame() *game {
	g := &game{}
	g.Id = Uuid()
	g.State = "new"
	return g
}

func (g *game) getState() []byte {
	gs := &gameStateMsg{}
	gs.Type = "gamestate"
	gs.Id = g.Id
	gs.StartTime = g.StartTime
	gs.TimeLeftMs = g.StartTime.Add(g.TimeLimit).Sub(g.LastUpdate).Seconds()
	gs.Players = g.Players
	gs.Bullets = g.Bullets
	gs.Tiles = g.Tiles
	gs.State = g.State
	b, err := json.Marshal(gs)
	if err != nil {
		log.Println("error marshaling:" + err.Error())
	}
	return b
}

func (g *game) getStateForPlayer(p *player) []byte {
	// TODO: Hide occluded players from gamestate sent to AI. (not critical)
	return g.getState()
}

func (g *game) start() error {
	if g.State == "new" {
		log.Println("game starting")

		// TODO: Randomize player positions!

		g.State = "running"
		g.StartTime = time.Now()
		return nil
	} else {
		return errors.New("Can't start game, state is " + g.State)
	}
}

func (g *game) end() error {
	log.Println("game ending")
	g.State = "end"
	return nil
}

func (g *game) hasEnded() bool {
	if time.Since(g.StartTime) > g.TimeLimit {
		log.Println("Outtatime")
		return true
	}
	alivePlayers := 0
	for _, k := range g.Players {
		if k.Hitpoints > 0 {
			alivePlayers++
		}
	}
	if alivePlayers < 2 {
		log.Println("Outtaplayers, ending game.")
		//return true
	}
	return false
}

func (g *game) newPlayer(position *Vector, name string) (*player, error) {
	if g.State != "new" {
		return nil, errors.New("Can't create new player when game state is " + g.State)
	}
	p := &player{}
	p.Id = Uuid()
	p.Hitpoints = 100
	p.Radius = 1
	p.Position = position
	p.Velocity = &Vector{0, 0}
	p.Name = name
	g.Players = append(g.Players, p)
	p.LookingAt = &Vector{p.Position.X, p.Position.Y}
	log.Println("newplayer", p.Position, p.Velocity)
	return p, nil
}

func (g *game) newBullet(position *Vector, velocity *Vector, firedBy string) *bullet {
	b := &bullet{}
	b.Id = Uuid()
	b.Damage = 10
	b.Radius = 0.1
	b.Position = position
	b.Velocity = velocity
	b.FiredBy = firedBy
	g.Bullets = append(g.Bullets, b)
	return b
}

func (g *game) newTile(position *Vector, width, height float64) *tile {
	t := &tile{}
	t.Id = Uuid()
	t.Width = width
	t.Height = height
	t.Position = position
	t.Velocity = &Vector{0, 0}
	g.Tiles = append(g.Tiles, t)
	return t
}

func (g *game) rmPlayer(p *player) {
	for i, k := range g.Players {
		if k == p {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
		}
	}
}

func (g *game) rmBullet(b *bullet) {
	for i, k := range g.Bullets {
		if k == b {
			g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
		}
	}
}

func (g *game) rmTile(t *tile) {
	for i, k := range g.Tiles {
		if k == t {
			g.Tiles = append(g.Tiles[:i], g.Tiles[i+1:]...)
		}
	}
}

func (g *game) update(dt time.Duration) {
	switch g.State {
	case "new":
	case "running":
		g.Collisions = nil
		for _, b := range g.Bullets {
			b.update(g, dt)
			// Remove all bullets that are dead.
			var newBullets []*bullet
			for _, k := range g.Bullets {
				if !k.Dead {
					newBullets = append(newBullets, k)
				}
			}
			g.Bullets = newBullets
		}
		for _, p := range g.Players {
			p.update(g, dt)
		}

		g.LastUpdate = time.Now()
		if g.hasEnded() {
			g.end()
		}
	case "ended":
	}
}
