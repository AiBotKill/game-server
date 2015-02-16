package main

import (
	"log"
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

	natsInit()

	// Logging server stats until the server is stopped.
	for {
		<-time.After(time.Second * 1)
		err := natsEncodedConn.Publish("ping", map[string]interface{}{
			"ping":      "gameServer",
			"serviceId": serviceId,
			"time":      time.Now().String(),
		})
		if err != nil {
			log.Println(err.Error())
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

// natsInit inits the root addresses that are needed to communicate to gameserver.
func natsInit() {
	err := natsEncodedConn.Publish("newGameServer", map[string]interface{}{
		"id": serviceId,
	})
	if err != nil {
		log.Panicln("Can't publish newGameserver message:", err.Error())
	}

	// Subscribe to createGame
	natsEncodedConn.Subscribe("createGame", func(subj string, reply string, msg *CreateGameMsg) {
		// Create game
		g := newGame()
		g.GameArea = msg.GameArea
		g.TimeLimit = time.Duration(msg.TimeLimit) * time.Second
		g.StartingPositions = msg.StartingPositions
		g.Mode = msg.Mode

		for _, t := range msg.Tiles {
			pos := &Vector{t.X + 0.5, t.Y + 0.5}
			g.newTile(pos, 1, 1)
		}

		// Subscribe to gameId.join
		natsEncodedConn.Subscribe(g.Id+".join", func(subj string, reply string, msg *JoinMsg) {
			log.Println("join request from", msg.BotId)
			p, err := g.newPlayer(&Vector{0, 0}, msg.Name)
			if err != nil {
				log.Println("Error creating player:", err.Error())
			}
			p.BotId = msg.BotId
			log.Println("join from", msg)
			natsEncodedConn.Publish(reply, Reply{Status: "ok", Id: g.Id})
		})

		// Subscribe to gameId.start
		natsEncodedConn.Subscribe(g.Id+".start", func(subj string, reply string, msg *StartGameMsg) {
			err := g.start()
			if err != nil {
				natsEncodedConn.Publish(reply, Reply{Status: "error", Error: err.Error()})
			}
			natsEncodedConn.Publish(reply, Reply{Status: "ok", Id: g.Id})
		})

		// Subscribe to gameId.start
		natsEncodedConn.Subscribe(g.Id+".end", func(subj string, reply string, msg *EndGameRequest) {
			err := g.end()
			if err != nil {
				natsEncodedConn.Publish(reply, Reply{Status: "error", Error: err.Error()})
			}
			natsEncodedConn.Publish(reply, Reply{Status: "ok", Id: g.Id})
		})

		// Invite players
		for _, p := range msg.Players {
			natsEncodedConn.Publish("aiserver."+p.BotId, JoinRequest{
				Type:     "joinRequest",
				GameId:   g.Id,
				GameMode: msg.Mode,
			})
		}

		// Reply game creator with id.
		natsEncodedConn.Publish(reply, &Reply{
			Id:     g.Id,
			Status: "ok",
		})

		go func() {
			for {
				<-time.After(time.Second)
				g.update()
				log.Println("game update: " + g.State)
				err := natsEncodedConn.Publish(g.Id+".gamestate", string(g.getState()))
				if err != nil {
					log.Println(err.Error())
				}

				if g.State == "end" {
					log.Println("Game " + g.Id + " has ended.")
					natsEncodedConn.Publish(g.Id+"gameEnd", string(g.getState()))
					return
				}
			}
		}()
	})
}

type GameStateMsg struct {
	Id        string    `json:"id"`
	StartTime time.Time `json:"startTime"`
}

type JoinMsg struct {
	BotId string `json:"botId"`
	Name  string `json:"name"`
}

type Reply struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}

type JoinRequest struct {
	Type     string `json:"type"`
	GameId   string `json:"gameId"`
	GameMode string `json:"gameMode"`
}

// CreateGameMsg is sent by game-console to create a new game.
// Replied with IdReplyMsg
type CreateGameMsg struct {
	TimeLimit   int64      `json:"timeLimit"`
	Environment string     `json:"environment"`
	GameArea    [2]float64 `json:"gameArea"`
	Tiles       []struct {
		Type string  `json:"type"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
	} `json:"tiles"`
	StartingPositions []*Vector              `json:"startingPositions"`
	Mode              string                 `json:"mode"`
	Players           []*CreateGameMsgPlayer `json:"players"`
}

type CreateGameMsgPlayer struct {
	Team  int    `json:"team"`
	BotId string `json:"botId"`
}

// StartGameMsg is sent to "<gameId>.start", which will trigger start if the
// game state is not "new"
type StartGameMsg struct {
}

type EndGameRequest struct {
}
