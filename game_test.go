package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestGame(t *testing.T) {
	done := make(chan bool)
	if err := startNats(); err != nil {
		t.Log("Can connect or start to gnatsd:", err.Error())
	}
	//serviceId = Uuid()
	//natsInit()
	log.Println("natsinit done")

	err := natsEncodedConn.Publish("test", "test")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	natsEncodedConn.Subscribe("aiserver.B123", func(subj string, reply string, msg *JoinRequest) {
		log.Println("Got joinrequest to:", msg.GameId)
		<-time.After(time.Second * 1)

		var joinreply Reply
		joinMsg := JoinMsg{BotId: "B123", Name: "jack"}
		err := natsEncodedConn.Request(msg.GameId+".join", joinMsg, &joinreply, time.Second*1)

		if err != nil {
			log.Println("Err joinin:", err.Error())
			t.Fail()
			done <- true
			return
		}
		log.Println("Joinreply", joinreply)
		done <- true
	})

	natsEncodedConn.Subscribe("aiserver.B124", func(subj string, reply string, msg *JoinRequest) {
		log.Println("Got joinrequest to:", msg.GameId)
		<-time.After(time.Second * 1)

		var joinreply Reply
		joinMsg := JoinMsg{BotId: "B124", Name: "jack"}
		err := natsEncodedConn.Request(msg.GameId+".join", joinMsg, &joinreply, time.Second*1)

		if err != nil {
			log.Println("Err joinin:", err.Error())
			t.Fail()
			done <- true
			return
		}
		log.Println("Joinreply", joinreply)
		done <- true
	})

	req := &CreateGameMsg{TimeLimit: 60}
	var plrs []*CreateGameMsgPlayer
	plrs = append(plrs, &CreateGameMsgPlayer{BotId: "B123"})
	plrs = append(plrs, &CreateGameMsgPlayer{BotId: "B124"})
	req.Players = plrs

	var reply Reply
	err = natsEncodedConn.Request("createGame", req, &reply, time.Second*10)
	if err != nil || reply.Status != "ok" {
		t.Log(err.Error())
		t.Fail()
	}

	gameId := reply.Id
	log.Println("Game created:", gameId)

	b, err := json.Marshal(&reply)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	log.Println(string(b))

	<-done
	log.Println("done")

	var startReply Reply
	err = natsEncodedConn.Request(gameId+".start", StartGameMsg{}, &startReply, time.Second*4)
	if err != nil {
		log.Println(err.Error())
	}
	if startReply.Status != "ok" {
		log.Println(startReply.Error)
		t.Fail()
	}
	log.Println("start:", startReply.Status)

	<-time.After(time.Second * 30)
	var endReply Reply
	err = natsEncodedConn.Request(gameId+".end", EndGameRequest{}, &endReply, time.Second*4)
	if err != nil {
		log.Println(err.Error())
	}
	if endReply.Status != "ok" {
		log.Println(endReply.Error)
		t.Fail()
	}
	t.Log("end:", endReply.Status)
}
