package server

import (
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/player"
	. "github.com/olsdavis/goelan/protocol"
)

// This file contains all the handlers for the play state.

// clientSettingsHandler updates clients' settings.
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
	id := packet.ReadLong()
	sender.PendingKeepAlives.QueryAndComplete(func(test interface{}) bool {
		data := test.(player.KeepAliveData)
		return data.ID == id
	})
}

func chatMessageHandler(packet *RawPacket, sender *Connection) {
	message := sender.Player.Name + " > " + packet.ReadString()
	log.Info(sender.Player.Name, message)
	sender.SendMessage(message, ChatMessageMode)
}

func teleportConfirmHandler(packet *RawPacket, sender *Connection) {
	id := packet.ReadVarint()
	sender.PendingTeleportConfirmations.QueryAndComplete(func(test interface{}) bool {
		data := test.(player.TeleportConfirmData)
		return data.ID == id
	})
}

func playerPositionAndLookHandler(packet *RawPacket, sender *Connection) {
	//TODO: implement
}

func animationHandler(packet *RawPacket, sender *Connection) {
	//TODO: implement
}

func clickWindowHandler(packet *RawPacket, sender *Connection) {
	//TODO: implement
}

func closeWindowHandler(packet *RawPacket, sender *Connection) {
	//TODO: implement
}
