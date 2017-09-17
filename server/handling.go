package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/olsdavis/goelan/auth"
	"github.com/olsdavis/goelan/encrypt"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/player"
	. "github.com/olsdavis/goelan/protocol"
	"github.com/olsdavis/goelan/util"
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
			PluginMessagePacketId:     pluginMessageHandler,
			KeepAliveIncomingPacketId: keepAliveHandler,
			ClientSettingsPacketId:    clientSettingsHandler,
			ClientStatusPacketId:      clientStatusHandler,
			IncomingChatPacketId:      chatMessageHandler,
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

/*** HANDSHAKE HANDLERS ***/

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

/*** LOGIN HANDLERS ***/

// Handles the login start packet.
func loginStartHandler(packet *RawPacket, sender *Connection) {
	username := packet.ReadStringMax(16)
	response := NewResponse()

	// old client
	if sender.GetServer().GetServerVersion().ProtocolVersion > sender.ProtocolVersion {
		sender.Disconnect(fmt.Sprintf("Your client is outdated. I'm on %v.", sender.GetServer().GetServerVersion().Name))
		return
	} else if sender.GetServer().GetServerVersion().ProtocolVersion < sender.ProtocolVersion {
		sender.Disconnect(fmt.Sprintf("I'm still on %v.", sender.GetServer().GetServerVersion().Name))
		return
	}

	if sender.GetServer().IsOnlineMode() {
		// send encryption request
		token := encrypt.GenerateVerifyToken()
		publicKey := sender.GetServer().GetPublicKey()

		response.WriteString("")
		response.WriteByteArray(publicKey)
		response.WriteByteArray(token)
		sender.Write(response.ToRawPacket(EncryptionRequestPacketId))
		sender.VerifyToken = token
		sender.VerifyUsername = username
	} else {
	}
}

// Handles the encryption request packet.
func encryptionResponseHandler(packet *RawPacket, sender *Connection) {
	sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, sender.GetServer().GetPrivateKey(), packet.ReadByteArray())
	if err != nil {
		panic(err)
	}
	verifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, sender.GetServer().GetPrivateKey(), packet.ReadByteArray())
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(verifyToken, sender.VerifyToken) {
		sender.Disconnect("Invalid verify token.")
		return
	}

	// assign new readers and writers
	aesCipher, err := aes.NewCipher(sharedSecret)
	if err != nil {
		panic(err)
	}

	sender.Writer = cipher.StreamWriter{
		W: sender.Writer,
		S: encrypt.NewCFB8Encrypt(aesCipher, sharedSecret),
	}

	reader := cipher.StreamReader{
		R: sender.Reader.R,
		S: encrypt.NewCFB8Decrypt(aesCipher, sharedSecret),
	}

	sender.Reader.R = reader
	// end

	// auth
	profile, err := auth.Auth(sender.VerifyUsername, sharedSecret, sender.GetServer().GetPublicKey())
	if err != nil {
		sender.Disconnect("Could not connect to Mojang servers.")
		log.Error("Error while connecting to Mojang servers:", err)
		return
	}

	// Login Success packet
	response := NewResponse()
	response.WriteString(util.ToHypenUUID(profile.UUID))
	response.WriteString(profile.Name)
	sender.Write(response.ToRawPacket(LoginSuccessPacketId))

	sender.SharedSecret = sharedSecret
	// release the data we don't need anymore
	sender.VerifyToken = emptyArray
	sender.VerifyUsername = ""

	if ok, reason := sender.GetServer().CanConnect(profile.Name, profile.UUID); !ok {
		sender.Disconnect(reason)
		return
	}

	// New connection state
	sender.ConnectionState = PlayState
	AssignHandler(sender)
	response.Clear()
	// Join Game packet
	response.WriteInt(0). // entity id
				WriteUnsignedByte(1).   // gamemode
				WriteInt(0).            // dimension
				WriteUnsignedByte(0).   // difficulty
				WriteUnsignedByte(0).   // max players (ignored)
				WriteString("default"). // level type
				WriteBoolean(false)     // reduced debug info
	sender.Write(response.ToRawPacket(JoinGamePacketId))
	response.Clear()
	sender.GetServer().FinishLogin(*profile, sender)
}

/*** PLAY HANDLERS ***/

func clientSettingsHandler(packet *RawPacket, sender *Connection) {
	sender.Player.Settings.Locale = packet.ReadStringMax(16)
	sender.Player.Settings.ViewDistance = packet.ReadByte()
	sender.Player.Settings.ChatMode = player.ChatMode(packet.ReadVarint())
	sender.Player.Settings.ColorsEnabled = packet.ReadBoolean()
	sender.Player.Settings.DisplayedSkinParts = packet.ReadUnsignedByte()
	sender.Player.Settings.MainHand = player.Hand(packet.ReadVarint())
}

func clientStatusHandler(packet *RawPacket, sender *Connection) {
}

func pluginMessageHandler(packet *RawPacket, sender *Connection) {
}

func keepAliveHandler(packet *RawPacket, sender *Connection) {
	id := packet.ReadVarint()
	if sender.LastKeepAlive.ID == id {
		sender.Lock()
		sender.LastKeepAlive.ID = -1
		sender.Unlock()
	}
}

func chatMessageHandler(packet *RawPacket, sender *Connection) {
	message := sender.Player.Name + " > " + packet.ReadString()
	log.Info(sender.Player.Name, message)
	sender.SendMessage(message, ChatMessageMode)
}
