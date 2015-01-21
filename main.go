package main

import (
	"log"
	"runtime"
	"time"
)

var serviceId string

func main() {
	if err := startNats(); err != nil {
		log.Panicln("Can connect or start to gnatsd:", err.Error())
	}
	id := Uuid()

	err := natsEncodedConn.Publish("new_gameserver", map[string]interface{}{
		"gameserver_id": id,
	})

	if err != nil {
		log.Panicln("Can't publish new_gameserver message.")
	}

	natsEncodedConn.QueueSubscribe("createGame", "createGame", natsCreateGame)

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

func natsCreateGame(subj string, reply string, msg interface{}) {
	g := NewGame()

	log.Println(msg)

	natsEncodedConn.Subscribe(g.Id+".create_player", func(subj string, reply string, msg interface{}) {
		if g.State != "new" {
			natsEncodedConn.Publish(reply, map[string]interface{}{
				"error": "can't create player when gamestate is '" + g.State + "'.",
			})
		}

		p := NewEntity([2]float64{0, 0}, 1, g)
		natsEncodedConn.Subscribe(p.Id+".action", func(subj string, reply string, msg interface{}) {
			log.Println("player", p.Id, "actionMsg", msg)
		})

		natsEncodedConn.Publish(reply, p.Id)
	})

	natsEncodedConn.Subscribe(g.Id+".join", func(subj string, reply string, msg interface{}) {})

	natsEncodedConn.Publish(reply, g.Id)

	err := natsEncodedConn.Publish("new_game", map[string]interface{}{
		"gameId": g.Id,
	})

	if err != nil {
		log.Println("ERROR: while publishing new_game message:", err.Error())
	}

}
