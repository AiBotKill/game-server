package main

import (
	"log"
	"time"
)

type game struct {
	LastUpdate time.Time
	Entities   []*entity
}

func (g *game) update(dt time.Duration) {
	//dt := time.Since(g.LastUpdate)
	g.LastUpdate = time.Now()
	for _, e := range g.Entities {
		e.update(dt)
	}
}

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
	Entity *entity
	Point  [2]float64
}

func (g *game) collision(vector [4]float64) []*collision {
	// Loop trough all entities, optimize opportunity with culling.
	var collisions []*collision
	for _, e := range g.Entities {
		xmin := e.Location[0] - (e.Dimensions[0] / 2.0)
		ymin := e.Location[1] - (e.Dimensions[1] / 2.0)
		xmax := e.Location[0] + (e.Dimensions[0] / 2.0)
		ymax := e.Location[1] + (e.Dimensions[1] / 2.0)

		borders := [][4]float64{[4]float64{xmin, ymin, xmax, ymin},
			[4]float64{xmin, ymin, xmin, ymax},
			[4]float64{xmax, ymin, xmax, ymax},
			[4]float64{xmin, ymax, xmax, ymax}}

		for _, b := range borders {
			if i := intersectionPoint(vector[0], vector[1], vector[2], vector[3], b[0], b[1], b[2], b[3]); i != nil {
				log.Println(b)
				collisions = append(collisions, &collision{e, [2]float64{i[0], i[1]}})
			}
		}
	}
	return collisions
}
