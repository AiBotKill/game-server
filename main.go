package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/apcera/nats"
)

const (
	AI_TIMEOUT = time.Second * 2
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
		<-time.After(time.Second * 5)
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
	natsConn.Subscribe("startGame", func(msg *nats.Msg) {
		var subs []*nats.Subscription

		var createGameMsg CreateGameMsg
		err := json.Unmarshal(msg.Data, &createGameMsg)
		if err != nil {
			natsConn.Publish(msg.Reply, NewReply("", err))
			return
		}

		// Create game
		g := newGame()
		g.GameArea = createGameMsg.GameArea
		g.TimeLimit = time.Duration(createGameMsg.TimeLimit) * time.Second
		g.StartingPositions = createGameMsg.StartingPositions
		g.Mode = createGameMsg.Mode

		for _, t := range createGameMsg.Tiles {
			pos := &Vector{t.X + 0.5, t.Y + 0.5}
			g.newTile(pos, 1, 1)
		}

		for _, jp := range createGameMsg.Players {
			p, err := g.newPlayer(&Vector{0, 0}, "")
			if err != nil {
				log.Println("ERROR:", err.Error())
				natsConn.Publish(msg.Reply, NewReply(g.Id, err))
				return
			}

			p.BotId = jp.BotId
			p.Team = jp.Team
			rnd := rand.Intn(len(createGameMsg.StartingPositions))
			p.Position = createGameMsg.StartingPositions[rnd]
			createGameMsg.StartingPositions = append(createGameMsg.StartingPositions[:rnd], createGameMsg.StartingPositions[rnd+1:]...)

			// Subscribe to "botId.action" address.
			if sub, err := natsConn.Subscribe(p.BotId+".action", func(msg *nats.Msg) {
				var action ActionMsg
				if err := json.Unmarshal(msg.Data, &action); err != nil {
					log.Println("ERROR:", err.Error())
					natsConn.Publish(msg.Reply, NewReply(g.Id, err))
					return
				}
				p.Action.Type = action.Type
				p.Action.Direction = action.Direction
				natsConn.Publish(msg.Reply, NewReply(g.Id, nil))
			}); err != nil {
				natsConn.Publish(msg.Reply, NewReply(g.Id, err))
				log.Println("ERROR:", err.Error())
				return
			} else {
				subs = append(subs, sub)
			}
		}

		// Subscribe to gameId.end
		if sub, err := natsConn.Subscribe(g.Id+".end", func(msg *nats.Msg) {
			err := g.end()
			natsConn.Publish(msg.Reply, NewReply(g.Id, err))
		}); err != nil {
			natsConn.Publish(msg.Reply, NewReply(g.Id, err))
			log.Println(err.Error())
			return
		} else {
			subs = append(subs, sub)
		}

		// Start game and return message.
		if err := g.start(); err != nil {
			natsConn.Publish(msg.Reply, NewReply(g.Id, err))
			return
		} else {
			natsConn.Publish(msg.Reply, NewReply(g.Id, err))
		}

		// Give some time to open visualization before starting for real.
		<-time.After(time.Second * 5)

		go func() {
			for {
				<-time.After(time.Second)
				g.update(time.Second) // TODO some logic for this!
				log.Println("game update: " + g.State)

				b := g.getState()
				if err := natsConn.Publish(g.Id+".gameState", b); err != nil {
					log.Println(string(b))
					log.Println("gamestate pub error: " + err.Error())
				}

				wg := &sync.WaitGroup{}
				for _, p := range g.Players {
					// Backpressure
					wg.Add(1)
					go func() {
						defer wg.Done()
						plrGamestate := g.getStateForPlayer(p)
						b, err := natsConn.Request(p.BotId+".gameState", plrGamestate, AI_TIMEOUT)
						if err != nil {
							log.Println("AI gamestate req error:", err.Error())
							p.Linkdead = true
							return
						} else {
							p.Linkdead = false
						}

						var reply Reply
						err = json.Unmarshal(b.Data, &reply)
						if err != nil {
							log.Println(err.Error())
							p.Linkdead = true
							return
						} else if reply.Error != "" {
							log.Println(reply.Error)
							p.Linkdead = true
						}
					}()
					wg.Wait()
				}

				if g.State == "end" {
					// Game has ended, clean up and publish gameEnd message.
					log.Println("Game " + g.Id + " has ended.")
					natsConn.Publish(g.Id+".gameEnd", g.getState())

					// Unsubscribe all subscriptions made during this game.
					for _, sub := range subs {
						if err := sub.Unsubscribe(); err != nil {
							log.Println(err.Error())
						}
					}
					return
				}
			}
		}()
	})
}

func NewReply(id string, err error) []byte {
	reply := &Reply{
		Type: "reply",
		Id:   id,
	}
	if err != nil {
		reply.Status = "error"
		reply.Error = err.Error()
	} else {
		reply.Status = "ok"
	}

	b, _ := json.Marshal(&reply)
	return b
}

type JoinMsg struct {
	BotId string `json:"botId"`
	Name  string `json:"name"`
}

type Reply struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}

type JoinRequest struct {
	Type     string `json:"type"`
	GameId   string `json:"gameId"`
	GameMode string `json:"gameMode"`
}

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
	Team  int64  `json:"team"`
	BotId string `json:"botId"`
}
