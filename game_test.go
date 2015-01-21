package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestEntityMovement(t *testing.T) {
	Convey("Given a new entity with acceleration of [1, 0.5]", t, func() {
		e := &entity{}
		e.Dimensions = [2]float64{1, 1}
		e.Location = [2]float64{0, 0}
		e.Acceleration = [2]float64{1, 0.5}
		e.Velocity = [2]float64{0, 0}
		Convey("after one second", func() {
			e.update(time.Second * 1.0)

			Convey("x velocity should be 1", func() {
				So(e.Velocity[0], ShouldEqual, 1.0)
			})

			Convey("y velocity should be 0.5", func() {
				So(e.Velocity[1], ShouldEqual, 0.5)
			})

			t.Log("location", e.Location)
		})

		Convey("after two seconds", func() {
			e.update(time.Second * 2.0)

			Convey("x velocity should be 2", func() {
				So(e.Velocity[0], ShouldEqual, 2.0)
			})

			Convey("y velocity should be 1.0", func() {
				So(e.Velocity[1], ShouldEqual, 1.0)
			})

			Convey("X location should be 4", func() {
				So(e.Location[0], ShouldEqual, 4.0)
			})

			Convey("Y location should be 2", func() {
				So(e.Location[1], ShouldEqual, 2.0)
			})
		})

		Convey("Action", func() {
			e.action("tooshort")
			e.action("look 1.0 1.0")
			e.action("shoot 1.0 1.0")
			e.action("shoot 1.0 xxxx")
			e.action("shoot xxxx 1.0")

			e.action("move -1 0")
			e.action("move 0 -1")

			e.action("move -1 -1")
			e.action("move 0 0")
			e.action("move 1 1")

			e.action("move 0 1")
			e.action("move 1 0")

			e.action("move -1 1")
			e.action("move 1 -1")

		})

	})
}

func TestGameIntersection(t *testing.T) {
	Convey("Given a new game with one entity moving trough a line on X axis", t, func() {
		game := &game{}
		e := &entity{}
		e.Dimensions = [2]float64{1, 1}
		e.Location = [2]float64{1, 1}
		e.Game = game

		game.Entities = append(game.Entities, e)
		line := [4]float64{0, 0, 2, 2}

		Convey("intersection should happen", func() {
			collisions := game.collision(line)
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
