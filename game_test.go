package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestPlayerShooting(t *testing.T) {
	Convey("Given a new game with two players", t, func() {
		g := &game{}
		g.State = "running"
		p1 := g.newPlayer([2]float64{0, 0}, "player 1")
		_ = g.newPlayer([2]float64{10, 0}, "player 2")
		p1.shoot(Vector{10, 0})
		g.update(time.Second * 10)
	})
}

func TestEntityMovement(t *testing.T) {
	Convey("Given a new player with Velocity of [1, 0.5]", t, func() {
		e := &player{}
		e.Dimensions = [2]float64{1, 1}
		e.Location = [2]float64{0, 0}
		e.Velocity = [2]float64{1, 0.5}
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

			Convey("X location should be 2", func() {
				So(e.Location[0], ShouldEqual, 2.0)
			})

			Convey("Y location should be 1", func() {
				So(e.Location[1], ShouldEqual, 1.0)
			})
		})

	})
}

func TestGameIntersection(t *testing.T) {
	Convey("Given a new game with one player moving trough a line on X axis", t, func() {
		game := &game{}
		game.State = "running"
		e := &player{}
		e.Dimensions = [2]float64{1, 1}
		e.Location = [2]float64{1, 1}
		e.Game = game

		game.Players = append(game.Players, e)
		line := [4]float64{0, 0, 2, 2}

		Convey("intersection should happen", func() {
			collisions := game.collision(line)
			So(len(collisions), ShouldBeGreaterThan, 0)
		})

	})

	Convey("Given a new game with one player moving trough a line on X axis", t, func() {
		game := &game{}
		game.State = "running"
		e := &player{}
		e.Dimensions = [2]float64{1.0, 1.0}
		e.Location = [2]float64{0.0, 0.0}
		e.Velocity = [2]float64{1.0, 0.0}
		game.Players = append(game.Players, e)
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

	Convey("Given a new game with one player moving trough a line on Y axis", t, func() {
		game := &game{}
		game.State = "running"
		e := &player{}
		e.Dimensions = [2]float64{1.0, 1.0}
		e.Location = [2]float64{0.0, 0.0}
		e.Velocity = [2]float64{0.1, 0.9}
		game.Players = append(game.Players, e)
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
