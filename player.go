package main

import (
	"encoding/json"
	"log"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
)

type Player struct {
	Name          string
	Uuid          uuid.UUID
	Conn          *net.Conn
	Upstream      *UpstreamConn
	CurrentServer string
}

func (pl *Player) HardDisconnect() {
	pl.Conn.Close()
	pl.Upstream.ClientConn.Close()
}

func (pl *Player) DisconnectWithMessage(msg chat.Message) {
	m, _ := json.Marshal(msg)
	pl.Conn.WritePacket(pk.Marshal(
		0x1a,
		pk.Chat(string(m)),
	))
	pl.HardDisconnect()
}

func (pl *Player) Handle() {
	go func() {
		for {
			p := pl.Upstream.ReadPacket()
			if p == nil {
				break
			}
			pl.Conn.WritePacket(*p)
		}
	}()

	for {
		var p pk.Packet
		if err := pl.Conn.ReadPacket(&p); err != nil {
			log.Printf("ReadPacket error: %v", err)
			pl.Conn.Close()
			pl.Upstream.ClientConn.Close()
			break
		}

		pl.Upstream.ClientConn.WritePacket(p)
	}
}
