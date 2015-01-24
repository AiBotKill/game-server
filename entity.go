package main

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func NewEntity(location [2]float64, entityType int, g *game) *entity {
	e := &entity{}
	e.Id = Uuid()
	e.Location = location
	e.Velocity = [2]float64{0, 0}
	e.Acceleration = [2]float64{0, 0}
	e.Dimensions = [2]float64{1, 1}
	e.EntityType = entityType
	e.Game = g
	e.HitPoints = 100.0
	return e
}

type entity struct {
	Id           string
	Game         *game
	HitPoints    float64
	Location     [2]float64
	Velocity     [2]float64
	Acceleration [2]float64
	Looking      [2]float64
	Dimensions   [2]float64 // [width, height float64]
	EntityType   int
	Exhausted    bool
}

func (e *entity) update(dt time.Duration) {
	for i := 0; i < len(e.Location); i++ {
		e.Velocity[i] += e.Acceleration[i] * dt.Seconds()
		e.Location[i] += e.Velocity[i] * dt.Seconds()
	}
	e.Exhausted = false
}

// Action is a string, separated with spaces in format of:
// MOVE  <X float64>,<Y float64>
// LOOK  <X float64>,<Y float64>
// SHOOT <X float64>,<Y float64>
func (e *entity) action(action string) error {
	tokens := strings.SplitN(action, " ", -1)
	if len(tokens) != 3 {
		return errors.New("Action should have 3 parts, separated with space")
	}

	cmd := tokens[0]
	x, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return err
	}

	y, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return err
	}

	a := math.Atan2(y, x) / math.Pi / 2 * 360

	d := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))

	switch tokens[0] {
	case "move":
		log.Println(cmd, x, y, a, d)
		break
	case "shoot":
		break
	case "look":
		break
	}
	e.Exhausted = true
	return nil
}
