package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGameIntersection(t *testing.T) {
	Convey("Given a new game with one entity moving trough a line on X axis", t, func() {
		game := &game{}
		e := &entity{}
		e.Dimensions = [2]float64{1, 1}
		e.Location = [2]float64{1, 1}

		game.Entities = append(game.Entities, e)
		line := [4]float64{0, 0, 2, 2}

		Convey("intersection should happen", func() {
			collisions := game.collision(line)
			t.Log("collisions:", collisions)
			So(len(collisions), ShouldBeGreaterThan, 0)
		})

	})

	Convey("Given a new game with one entity moving trough a line on X axis", t, func() {
		game := &game{}
		e := &entity{}
		e.Dimensions = [2]float64{1.0, 1.0}
		e.Location = [2]float64{0.0, 0.0}
		e.Velocity = [2]float64{1.0, 0.0}
		game.Entities = append(game.Entities, e)
		line := [4]float64{1.0, -1.0, 1.0, 1.0}

		Convey("immediately", func() {
			Convey("intersection shouldn't happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeZeroValue)
			})
		})

		Convey("after a second", func() {
			game.update(time.Second)
			Convey("intersection should happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeGreaterThan, 0)
			})
		})

		Convey("after a two seconds", func() {
			game.update(time.Second * 2)
			Convey("intersection should not happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeZeroValue)
			})
		})
	})

	Convey("Given a new game with one entity moving trough a line on Y axis", t, func() {
		game := &game{}
		e := &entity{}
		e.Dimensions = [2]float64{1.0, 1.0}
		e.Location = [2]float64{0.0, 0.0}
		e.Velocity = [2]float64{0.1, 0.9}
		game.Entities = append(game.Entities, e)
		line := [4]float64{-1.1, 1.0, 1.0, 1.0}

		Convey("immediately", func() {
			Convey("intersection shouldn't happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeZeroValue)
			})
		})

		Convey("after a second", func() {
			game.update(time.Second)
			Convey("intersection should happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeGreaterThan, 0)
			})
		})

		Convey("after a two seconds", func() {
			game.update(time.Second * 2)
			Convey("intersection should not happen", func() {
				collisions := game.collision(line)
				So(len(collisions), ShouldBeZeroValue)
			})
		})
	})

}
