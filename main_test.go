package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNatsConnectivity(t *testing.T) {
	Convey("Starting a nats server", t, func() {
		err := startNats()

		Convey("should not return an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("natsEncodedConn should NOT be nil", func() {
			So(natsEncodedConn, ShouldNotBeNil)
		})

		Convey("When subscribing to \"test_connectivity\"", func() {
			sub, err := natsEncodedConn.Subscribe("test_connectivity", func(subj string, reply string, msg string) {
				natsEncodedConn.Publish(reply, msg)
			})
			So(err, ShouldBeNil)

			Convey("It should reply when sent a message", func() {
				var response string
				err := natsEncodedConn.Request("test_connectivity", "hello test", &response, 1*time.Second)
				So(err, ShouldBeNil)
				So(response, ShouldEqual, "hello test")

			})
			So(sub.IsValid(), ShouldBeTrue)
			sub.Unsubscribe()
			So(sub.IsValid(), ShouldBeFalse)

		})

		Convey("init game nats subscriptions", func() {
			natsInit()
			var gameId string
			err := natsEncodedConn.Request("create_game", &GameMessage{}, &gameId, 1*time.Second)
			So(err, ShouldBeNil)

			Convey("It should reply whith gameid", func() {
				So(gameId, ShouldNotBeBlank)
				t.Log("game id:", gameId)
				Convey("It should reply whith gameid", func() {
					var playerId string
					err := natsEncodedConn.Request(gameId+".create_player", &GameMessage{}, &playerId, 1*time.Second)
					So(err, ShouldBeNil)
					t.Log("player id:", playerId)
				})
			})
		})

	})
}
