package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

var serviceId string

func main() {
	if err := startNats(); err != nil {
		log.Panicln("Can connect or start to gnatsd:", err.Error())
	}
	serviceId = Uuid()

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

func natsInit() {
	err := natsEncodedConn.Publish("new_gameserver", map[string]interface{}{
		"gameserver_id": serviceId,
	})

	if err != nil {
		log.Panicln("Can't publish new_gameserver message.")
	}

	natsEncodedConn.QueueSubscribe("create_game", "create_game", natsCreateGame)
	natsEncodedConn.Subscribe("new_gameserver", func(subj string, reply string, msg map[string]interface{}) {
		log.Println(msg)
	})
}

func natsCreateGame(subj string, reply string, msg GameMessage) {
	g := NewGame()
	natsEncodedConn.Subscribe(g.Id+".create_player", func(subj string, reply string, msg *GameMessage) {
		if g.State != "new" {
			natsEncodedConn.Publish(reply, map[string]interface{}{
				"error": "can't create player when gamestate is '" + g.State + "'.",
			})
		}

		// Crate new player
		p := g.newPlayer([2]float64{0, 0}, "jack")

		// Add action message handler to the new player
		natsEncodedConn.Subscribe(p.Id+".action", func(subj string, reply string, msg *GameMessage) {
			log.Println("player", p.Id, "actionMsg", msg)
			if msg.Action != nil {
				err := p.action(msg.Action.Type, Vector{msg.Action.Direction.X, msg.Action.Direction.Y})
				if err != nil {
					natsEncodedConn.Publish(reply, err)
					return
				}
				natsEncodedConn.Publish(reply, fmt.Sprintf("%s %f %f", msg.Action.Type, msg.Action.Direction.X, msg.Action.Direction.Y))
				return
			}
			natsEncodedConn.Publish(reply, "action was nil")
		})

		// Reply with player id
		natsEncodedConn.Publish(reply, p.Id)
	})

	// Subscribe game to "join" message
	natsEncodedConn.Subscribe(g.Id+".join", func(subj string, reply string, msg *GameMessage) {
		// TODO
	})

	// Subscribe game to "start" message
	natsEncodedConn.Subscribe(g.Id+".start", func(subj string, reply string, msg *GameMessage) {
		// TODO
	})

	// Reply with game id
	natsEncodedConn.Publish(reply, g.Id)

	// Publish "new game message"
	err := natsEncodedConn.Publish("new_game", map[string]interface{}{
		"gameId": g.Id,
	})

	if err != nil {
		log.Println("ERROR: while publishing new_game message:", err.Error())
	}

}
