package server

import (
	. "github.com/olsdavis/goelan/protocol"
	"github.com/olsdavis/goelan/log"
)

// This file contains all the handlers for the handshake state.

// Handles the handshake.
func handshakeHandler(packet *RawPacket, sender *Connection) {
	sender.ProtocolVersion = packet.ReadUnsignedVarint()
	// omit the following data, we don't need it
	packet.ReadStringMax(255)
	packet.ReadUnsignedShort()
	// end
	nextState := packet.ReadUnsignedVarint()
	switch nextState {
	// Status (server list)
	case HandshakeStatusNextState:
		response := NewResponse()
		version := sender.GetServer().GetServerVersion()
		list := ServerListPing{
			Ver: Version{Name: version.Name, Protocol: version.ProtocolVersion},
			Pl: Players{Max: sender.GetServer().GetMaxPlayers(),
				Online: sender.GetServer().GetOnlinePlayersCount()},
			Desc: ChatComponent{Text: sender.GetServer().GetMotd()},
			Fav:  "",
		}
		if sender.GetServer().HasFavicon() {
			list.Fav = sender.GetServer().GetFavicon()
		}
		response.WriteJSON(list)
		sender.Write(response.ToRawPacket(HandshakePacketId))
		// Login (wants to play)
	case HandshakeLoginNextState:
		sender.ConnectionState = LoginState
		AssignHandler(sender)
		// Unknown
	default:
		log.Error("Client handshake next state:", nextState)
	}
}

// Handles the ping packet. Sends back a pong packet with the received payload.
func pingPongHandler(packet *RawPacket, sender *Connection) {
	payload := packet.ReadLong()
	response := NewResponse()
	response.WriteLong(payload)
	sender.Write(response.ToRawPacket(PingPacketId))
}
