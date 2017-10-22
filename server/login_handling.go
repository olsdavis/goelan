package server

import (
	"fmt"
	"github.com/olsdavis/goelan/encrypt"
	"crypto/rsa"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/olsdavis/goelan/auth"
	. "github.com/olsdavis/goelan/protocol"
	"crypto/rand"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/util"
)

// This file contains all the handlers for the login state.

// Handles the login start packet.
func loginStartHandler(packet *RawPacket, sender *Connection) {
	username := packet.ReadStringMax(16)
	if sender.GetServer().GetServerVersion().ProtocolVersion > sender.ProtocolVersion {
		// old version
		sender.Disconnect(fmt.Sprintf("Your client is outdated. I'm on %v.", sender.GetServer().GetServerVersion().Name))
		return
	} else if sender.GetServer().GetServerVersion().ProtocolVersion < sender.ProtocolVersion {
		// new version
		sender.Disconnect(fmt.Sprintf("I'm still on %v.", sender.GetServer().GetServerVersion().Name))
		return
	}

	// send encryption request
	response := NewResponse()
	token := encrypt.GenerateVerifyToken()
	var encryptionRequest = struct {
		ServerID    string
		PublicKey   []byte
		VerifyToken []byte
	}{
		"",
		sender.GetServer().GetPublicKey(),
		token,
	}
	response.WriteStructure(encryptionRequest)
	sender.Write(response.ToRawPacket(EncryptionRequestPacketId))
	sender.VerifyToken = token
	sender.VerifyUsername = username
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
	// auth
	profile, err := auth.Auth(sender.VerifyUsername, sharedSecret, sender.GetServer().GetPublicKey())
	if err != nil {
		sender.Disconnect("Could not connect to Mojang servers.")
		log.Error("Error while connecting to Mojang servers:", err)
		return
	}
	// Login Success packet
	response := NewResponse()
	var loginSuccess = struct {
		UUID string
		Name string
	}{
		util.ToHypenUUID(profile.UUID),
		profile.Name,
	}
	response.WriteStructure(loginSuccess)
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
	var joinGame = struct {
		EntityId         int
		Gamemode         uint8
		Dimension        int
		Difficulty       uint8
		MaxPlayers       uint8
		LevelType        string
		ReducedDebugInfo bool
	}{
		0,
		0,
		0,
		0,
		0,
		"default",
		false,
	}
	sender.Write(response.WriteStructure(joinGame).ToRawPacket(JoinGamePacketId))
	response.Clear()
	sender.GetServer().FinishLogin(*profile, sender)
}
