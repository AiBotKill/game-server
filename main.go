package main

import (
	"log"
	"runtime"
	"time"
)

// servoiceId acts as the Id for this instance of game-server.
// It is used to communicate directly to this gameserver, using
// proper nats-address.
var serviceId string

func main() {
	if err := startNats(); err != nil {
		log.Panicln("Can connect or start to gnatsd:", err.Error())
	}
	serviceId = Uuid()

	// Logging server stats until the server is stopped.
	goRoutines := 0
	for {
		nim := natsEncodedConn.Conn.InMsgs
		nib := natsEncodedConn.Conn.InBytes
		nom := natsEncodedConn.Conn.OutMsgs
		nob := natsEncodedConn.Conn.OutBytes
		log.Println(nim, nib, nom, nob)

		time.Sleep(time.Second * 1)
		if goRoutines != runtime.NumGoroutine() {
			goRoutines = runtime.NumGoroutine()
			log.Println("Goroutines [", goRoutines, "]")
		}
	}
}

// IdReplyMsg acts as a reply to most requests, status is either "ok" or "error".
// If status is error, id is undefined.
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}

// RegisterAiMsg is sent by AI-Server when a new AI connects and registers.
// Replied with IdReplyMsg
type RegisterAiMsg struct {
	BotId string `json:"botId"`
}

// CreateGameMsg is sent by game-console to create a new game.
// Replied with IdReplyMsg
type CreateGameMsg struct {
	TimeLimit         time.Duration `json:"timelimit"`
	GameArea          [2]float64    `json:"gameArea"`
	Tiles             []*tile       `json:"tiles"`
	StartingPositions []*Vector     `json:"startingPositions"`
}

// CreatePlayerMsg is sent to "<gameId>.createPlayer" address, and creates a
// new player when the game is not already running or ended.
// Replied with IdReplyMsg
type CreatePlayerMsg struct {
	BotId string `json:"botId"`
	Name  string `json:"name"`
}

// ActionMsg is sent by AI to the "<playerId>.action" address, and commands the
// player entity what to do during next update. If no message is received before
// update, the last action is used.
// Replied with IdReplyMsg, with playerId as Id.
type ActionMsg struct {
	Type      string  `json:"type"`
	Direction *Vector `json:"direction"`
}

// StartGameMsg is sent to "<gameId>.start", which will trigger start if the
// game state is not "new"
type StartGameMsg struct {
}

// natsInit inits the root addresses that are needed to communicate to gameserver.
func natsInit() {

	// Publish a "newGameServer" message if someone listening cares (logging?)
	err := natsEncodedConn.Publish("newGameServer", map[string]interface{}{
		"id": serviceId,
	})
	if err != nil {
		log.Panicln("Can't publish newGameserver message:", err.Error())
	}

	// When AI-Server accepts an AI message, publish the news for console.
	natsEncodedConn.Subscribe("registerAi", func(subj string, reply string, msg *RegisterAiMsg) {
		natsEncodedConn.Publish("newGameserver", map[string]interface{}{
			"botId": msg.BotId,
		})
	})

	// When gameconsole requests a new game, reply with game id.
	natsEncodedConn.QueueSubscribe("createGame", "gameservers", func(subj string, reply string, msg *CreateGameMsg) {
		g := newGame(msg.TimeLimit, msg.GameArea)

		// Handler for "<gameId>.createPlayer" messages.
		_, err := natsEncodedConn.Subscribe(g.Id+".createPlayer", func(subj string, reply string, msg *CreatePlayerMsg) {

			// If the game state is not new, reply with error.
			if g.State != "new" {
				natsEncodedConn.Publish(reply, IdReplyMsg{
					Status: "error",
					Error:  "Can't create players, gamestate is: " + g.State,
					Id:     "",
				})
				return
			}

			// Create the player, location can be 0,0 as it will be randomized when
			// game is started.
			p, err := g.newPlayer(&Vector{0, 0}, msg.Name)
			if err != nil {
				natsEncodedConn.Publish(reply, IdReplyMsg{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			// Handler for "<playerId>.action" messages.
			_, err = natsEncodedConn.Subscribe(g.Id+".action", func(subj string, reply string, msg *ActionMsg) {
				p.Action.Type = msg.Type
				p.Action.Direction = msg.Direction

				// Action registered, reply with ok.
				natsEncodedConn.Publish(reply, IdReplyMsg{
					Status: "ok",
					Id:     p.Id,
				})
			})
			if err != nil {
				log.Println(err.Error())
				return
			}

			// All ok, reply createplayer message with "ok".
			log.Println("OK ")
			natsEncodedConn.Publish(reply, IdReplyMsg{
				Status: "ok",
				Id:     g.Id,
			})
		})
		if err != nil {
			log.Println(err.Error())
			return
		}

		_, err = natsEncodedConn.Subscribe(g.Id+".start", func(subj string, reply string, msg *StartGameMsg) {
			err := g.start()
			if err != nil {
				natsEncodedConn.Publish(reply, IdReplyMsg{
					Status: "error",
					Error:  err.Error(),
					Id:     g.Id,
				})
			} else {
				natsEncodedConn.Publish(reply, IdReplyMsg{
					Status: "ok",
					Id:     g.Id,
				})
			}
		})
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Start a routine to keep updating game and publishing gamestate until game status is "ended"
		go func() {
			for {
				<-time.After(time.Second)
				g.update()

				natsEncodedConn.Publish(g.Id+"gamestate", g.getState())

				for _, p := range g.Players {
					natsEncodedConn.Publish(p.Id+"gamestate", g.getStateForPlayer(p))
				}

				if g.State == "end" {
					log.Println("Game " + g.Id + " has ended.")
					return
				}
			}
		}()

	})
}
