package main

import (
	_ "embed"
	"log"

	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
	"github.com/mattn/go-colorable"
)

const ProtocolVersion = 756
const MaxPlayer = 200

// Packet IDs
const (
	PlayerPositionAndLookClientbound = 0x38
	JoinGame                         = 0x26
	KeepAlive                        = 0x21
)

func main() {
	log.SetOutput(colorable.NewColorableStdout())
	initDefaultServerListResponse()
	l, err := net.ListenMC(":25566")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	log.Println("Listening on :25566")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Accept error: %v", err)
		}
		go acceptConn(conn)
	}
}

func acceptConn(conn net.Conn) {
	defer conn.Close()
	// handshake
	protocol, intention, err := handshake(conn)
	if err != nil {
		log.Printf("Handshake error: %v", err)
		return
	}

	switch intention {
	default: //unknown error
		log.Printf("Unknown handshake intention: %v", intention)
	case 1: //for status
		acceptListPing(conn)
	case 2: //for login
		handlePlaying(conn, protocol)
	}
}

type PlayerInfo struct {
	Name    string
	UUID    uuid.UUID
	OPLevel int
}

//go:embed DimensionCodec.snbt
var dimensionCodecSNBT string

//go:embed Dimension.snbt
var dimensionSNBT string

func joinGame(conn net.Conn) error {
	return conn.WritePacket(pk.Marshal(JoinGame,
		pk.Int(0),          // EntityID
		pk.Boolean(false),  // Is hardcore
		pk.UnsignedByte(1), // Gamemode
		pk.Byte(1),         // Previous Gamemode
		pk.VarInt(1),       // World Count
		pk.Ary{Len: 1, Ary: []pk.Identifier{"world"}},      // World Names
		pk.NBT(nbt.StringifiedMessage(dimensionCodecSNBT)), // Dimension codec
		pk.NBT(nbt.StringifiedMessage(dimensionSNBT)),      // Dimension
		pk.Identifier("world"),                             // World Name
		pk.Long(0),                                         // Hashed Seed
		pk.VarInt(MaxPlayer),                               // Max Players
		pk.VarInt(15),                                      // View Distance
		pk.Boolean(false),                                  // Reduced Debug Info
		pk.Boolean(true),                                   // Enable respawn screen
		pk.Boolean(false),                                  // Is Debug
		pk.Boolean(true),                                   // Is Flat
	))
}
