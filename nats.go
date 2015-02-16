package main

import (
	"os"
	"log"
	gnatsd "github.com/apcera/gnatsd/server"
	"github.com/apcera/nats"
)

var natsConn *nats.Conn
var natsEncodedConn *nats.EncodedConn

func debugGnatsd() {
	opts := gnatsd.Options{}
	s := gnatsd.New(&opts)
	go s.Start()
}

func startNats() error {
	addr := os.Getenv("GNATSD_PORT_4222_TCP_ADDR")
	port := os.Getenv("GNATSD_PORT_4222_TCP_PORT")
	url := ("nats://" + addr + ":" + port)
	log.Println(url)
	var c *nats.Conn
	natsc, err := nats.Connect(url)

	if err != nil {
		debugGnatsd()
		natsc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			return err
		} else {
			c = natsc
		}
	} else {
		c = natsc
	}

	nc, err := nats.NewEncodedConn(c, "json")
	if err != nil {
		return err
	}
	natsConn = c
	natsEncodedConn = nc
	return nil
}
