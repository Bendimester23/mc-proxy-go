package main

import (
	"log"

	"github.com/Tnze/go-mc/net"
)

var playersByName = make(map[string]*Player)

func handlePlaying(conn net.Conn, protocol int32) {
	// login, get player info
	info, err := acceptLogin(conn)
	if err != nil {
		log.Print("Login failed")
		return
	}

	loginSuccess(conn, info.Name, info.UUID)

	upstream := ConnectToUpstreamServer("localhost", 25565, info.Name)

	pl := &Player{
		Name:          info.Name,
		Uuid:          info.UUID,
		Conn:          &conn,
		Upstream:      upstream,
		CurrentServer: "lobby",
	}

	playersByName[info.Name] = pl

	pl.Handle()
}
