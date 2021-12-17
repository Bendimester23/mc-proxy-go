package main

import (
	"fmt"
	"log"

	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
)

type UpstreamConn struct {
	ClientConn *net.Conn
	/* ServerConn *net.Conn */
	name string
}

func ConnectToUpstreamServer(host string, port int32, name string) *UpstreamConn {
	conn, err := net.DialMC(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	err = conn.WritePacket(pk.Marshal(
		0x00,
		pk.VarInt(ProtocolVersion), // Protocol version
		pk.String(host),            // Host
		pk.UnsignedShort(port),     // Port
		pk.Byte(2),
	))
	if err != nil {
		return nil
	}

	err = conn.WritePacket(pk.Marshal(
		packetid.LoginStart,
		pk.String(name),
	))
	if err != nil {
		return nil
	}

	for {
		var p pk.Packet

		if err := conn.ReadPacket(&p); err != nil {
			log.Printf("Error reading packet from server: %s\n", err.Error())
			return nil
		}

		log.Printf("Recieved packet %d from upstream server", p.ID)

		if p.ID == packetid.LoginSuccess {
			log.Println("Logged in on the upstream server")
			break
		}
	}

	return &UpstreamConn{
		name:       name,
		ClientConn: conn,
	}
}

func (u *UpstreamConn) ReadPacket() *pk.Packet {
	var p pk.Packet
	if err := u.ClientConn.ReadPacket(&p); err != nil {
		return nil
	}
	return &p
}
