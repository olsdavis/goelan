package server

import (
	"github.com/olsdavis/goelan/log"
	. "github.com/olsdavis/goelan/protocol"
)

type PacketHandler func(packet *RawPacket, sender *Connection)

var (
	handlers map[ConnectionState]stateHandler
)

type stateHandler interface {
	callHandler(packet *RawPacket, sender *Connection)
}

type stateMapHandler struct {
	handlers map[uint64]PacketHandler
}

func (handler stateMapHandler) callHandler(packet *RawPacket, sender *Connection) {
	h, ok := handler.handlers[packet.ID]
	if ok {
		h(packet, sender)
	} else {
		log.Debug("Unhandled ID:", packet.ID)
		log.Debug("Unhandled Data:", packet.Data.Buf)
	}
}

func init() {
	handlers = make(map[ConnectionState]stateHandler)

	// registers packet handlers
	handlers[HandshakeState] = stateMapHandler{
		map[uint64]PacketHandler{
			HandshakePacketId: handshakeHandler,
			PingPacketId:      pingPongHandler,
		},
	}
	handlers[LoginState] = stateMapHandler{
		map[uint64]PacketHandler{
			LoginStartPacketId:         loginStartHandler,
			EncryptionResponsePacketId: encryptionResponseHandler,
		},
	}
	handlers[PlayState] = stateMapHandler{
		map[uint64]PacketHandler{
			PluginMessagePacketId:                 pluginMessageHandler,
			KeepAliveIncomingPacketId:             keepAliveHandler,
			ClientSettingsPacketId:                clientSettingsHandler,
			ClientStatusPacketId:                  clientStatusHandler,
			IncomingChatPacketId:                  chatMessageHandler,
			TeleportConfirmPacketId:               teleportConfirmHandler,
			IncomingPlayerPositionAndLookPacketId: playerPositionAndLookHandler,
			IncomingAnimationPacketId:             animationHandler,
			ClickWindowPacketId:                   clickWindowHandler,
			CloseWindowPacketId:                   closeWindowHandler,
		},
	}
}

// AssignHandler assigns a new handler appropriated to the given
// connection's state.
func AssignHandler(conn *Connection) {
	h, ok := handlers[conn.ConnectionState]
	if ok {
		conn.Lock()
		conn.PacketHandler = h
		conn.Unlock()
	} else {
		log.Error("no handler found for state", conn.ConnectionState)
	}
}
