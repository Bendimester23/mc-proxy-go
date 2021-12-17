package main

import (
	"encoding/json"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
)

type ServerListResponse struct {
	Version     ServerListVersion `json:"version"`
	Players     ServerListPlayers `json:"players"`
	Description chat.Message      `json:"description"`
	FavIcon     string            `json:"favicon,omitempty"`
}

type ServerListVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerListPlayers struct {
	Max    int                      `json:"max"`
	Online int                      `json:"online"`
	Sample []ServerListPlayerEntity `json:"sample"`
}

type ServerListPlayerEntity struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

var defaultRespone = ""

var listPacket pk.Packet

func initDefaultServerListResponse() {

	a := ServerListResponse{
		Version: ServerListVersion{
			Name:     "SkyVillage Lobby",
			Protocol: ProtocolVersion,
		},
		Players: ServerListPlayers{
			Max:    MaxPlayer,
			Online: 0,
			Sample: []ServerListPlayerEntity{},
		},
		Description: chat.Message{
			Text:       "Epik SkyVillage szerver",
			Color:      "#c00f33",
			UnderLined: true,
		},
		FavIcon: "data:image/webp;base64,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",
	}
	b, _ := json.Marshal(a)
	defaultRespone = string(b)

	listPacket = pk.Marshal(0x00, pk.String(listResp()))
}

func acceptListPing(conn net.Conn) {
	var p pk.Packet
	for i := 0; i < 2; i++ {
		err := conn.ReadPacket(&p)
		if err != nil {
			return
		}

		switch p.ID {
		case 0x00:
			err = conn.WritePacket(listPacket)
		case 0x01:
			err = conn.WritePacket(p)
		}
		if err != nil {
			return
		}
	}
}

func listResp() string {
	return defaultRespone
}
