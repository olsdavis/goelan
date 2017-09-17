package server

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/olsdavis/goelan/encrypt"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/player"
	"github.com/olsdavis/goelan/protocol"
	"github.com/olsdavis/goelan/util"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"
	"math/rand"
	"github.com/olsdavis/goelan/world"
)

var (
	serverInstance *Server
)

// Returns the server's single instance
func Get() *Server {
	return serverInstance
}

const (
	banListFile    = "banlist.json"
	faviconFile    = "server-icon.png"
	propertiesFile = "server.toml"
)

// Server's properties, read from the properties file ("server.toml").
type ServerProperties struct {
	Port         uint16 `toml:"port"`        // server's port
	Address      string `toml:"address"`     // server's address
	Motd         string `toml:"motd"`        // server's motd (the description in the server list)
	MaxPlayers   int32  `toml:"max-players"` // the maximal amount of players that the server should host
	OnlineMode   bool   `toml:"online-mode"` // if true => authentication with Mojang servers
	ViewDistance int    `toml:"view-distance"`
}

// Represents a Minecraft server.
type Server struct {
	run         bool
	initialized bool             // true, if the server has been initialized
	properties  ServerProperties // server's properties

	clients    map[string]*Connection // online players
	playerLock sync.Mutex             // lock for the clients map

	BanList *player.BanList // contains the players that have been banned from the server

	serverVersion   ServerVersion   // server's version (protocol and name)
	favicon         string          // the favicon
	ticker          *time.Ticker    // the ticker for the ticking :)
	keepAliveTicker *time.Ticker    // the ticker used for sending keep alive packets
	rsaKeypair      *rsa.PrivateKey // the keypair used for encryption
	publicKey       []byte          // the public key in bytes

	world *world.World // one world only, for the moment

	ExitChan chan int // a channel used for server's close
}

// Creates a new server.
func CreateServer(properties ServerProperties) *Server {
	if serverInstance != nil {
		panic("already created a server")
	}
	serverInstance = &Server{
		run:             true,
		initialized:     false,
		properties:      properties,
		clients:         make(map[string]*Connection),
		playerLock:      sync.Mutex{},
		serverVersion:   ServerVersion{"1.12.1", 338},
		favicon:         "",
		ticker:          nil,
		keepAliveTicker: nil,
		rsaKeypair:      encrypt.GenerateKeyPair(),
		publicKey:       nil,
		world:           nil,
		ExitChan:        make(chan int, 1),
	}
	return serverInstance
}

// Creates a new server from the properties file.
func CreateServerFromProperties() *Server {
	props := readProperties()
	return CreateServer(*props)
}

func readProperties() *ServerProperties {
	var properties ServerProperties

	// properties file read
	if _, err := os.Open(propertiesFile); err != nil && os.IsNotExist(err) {
		log.Info(fmt.Sprintf("No %v file found. Creating one.", propertiesFile))

		properties = ServerProperties{
			Port:         25565,
			Address:      "127.0.0.1",
			Motd:         "A Goelan Minecraft server",
			MaxPlayers:   10,
			OnlineMode:   true,
			ViewDistance: 15,
		}

		f, e := os.Create(propertiesFile)
		if e != nil {
			log.Fatal(fmt.Sprintf("Could not create the '%v' file! %s", propertiesFile, e))
		}
		defer f.Close()
		toml.NewEncoder(f).Encode(properties)
	}

	if _, err := toml.DecodeFile(propertiesFile, &properties); err != nil {
		log.Error("Could not load configuration file 'server.properties'!", err)
		return nil
	}

	return &properties
}

// HasFavicon returns whether the server has a favicon image or not.
// (Which appears in the server list.)
func (s *Server) HasFavicon() bool {
	return s.favicon != ""
}

// GetFavicon returns the favicon; it may be empty - check before with HasFavicon().
func (s *Server) GetFavicon() string {
	return s.favicon
}

// GetMotd returns server's MOTD. (Which is the description in the server list.)
func (s *Server) GetMotd() string {
	return s.properties.Motd
}

// GetMaxPlayers returns the maximal amount of players the server can host.
// 0 if no limit. (There is no limit if max-players is set to 0 or less.)
func (s *Server) GetMaxPlayers() uint {
	if s.properties.MaxPlayers <= 0 {
		return 0
	}
	return uint(s.properties.MaxPlayers)
}

// IsOnlineMode returns whether the server must authenticate players with
// Mojang servers or not.
func (s *Server) IsOnlineMode() bool {
	return s.properties.OnlineMode
}

// GetPublicKey returns the public key. (Generates it if it has not been done yet.)
func (s *Server) GetPublicKey() []byte {
	if s.publicKey == nil {
		s.publicKey = encrypt.GeneratePublicKey(s.rsaKeypair)
	}
	return s.publicKey
}

// GetPrivateKey returns server's private key.
func (s *Server) GetPrivateKey() *rsa.PrivateKey {
	return s.rsaKeypair
}

// GetServerVersion returns server's version (protocol and name).
func (s *Server) GetServerVersion() ServerVersion {
	return s.serverVersion
}

// GetViewDistance returns server's view distance.
func (s *Server) GetViewDistance() int {
	return s.properties.ViewDistance
}

// Returns true if the server is currently running
func (s *Server) IsRunning() bool {
	return s.run
}

// Starts initializes the server.
func (s *Server) Start() {
	if s.initialized {
		return
	}

	log.Info(fmt.Sprintf("Protocol #%v (Minecraft %v)", s.serverVersion.ProtocolVersion, s.serverVersion.Name))

	// listen
	listen := fmt.Sprintf("%v:%v", s.properties.Address, s.properties.Port)
	socket, err := net.Listen("tcp", listen)

	if err != nil {
		panic(fmt.Sprintf("Could not create socket: %v", err))
	}

	s.load()

	s.initialized = true

	// 20 ticks per second
	s.ticker = time.NewTicker(time.Second / 20)
	s.keepAliveTicker = time.NewTicker(time.Second)
	go s.tick()
	go s.keepAlive()

	s.world = world.NewWorld("default")

	log.Info("Done start up! Waiting for players to join.")
	log.Info("Listening on", listen)
	for s.run {
		conn, _ := socket.Accept()
		go s.handleConnection(conn)
	}
}

// tick handles server's logic.
func (s *Server) tick() {
	for s.run {
		<-s.ticker.C
		// TODO: Logic
	}
}

// keepAlive handles the clients that should be kept
// alive or not.
func (s *Server) keepAlive() {
	for s.run {
		<-s.keepAliveTicker.C

		id := int32(rand.Intn(0xFFFE))
		packet := protocol.NewResponse().WriteVarint(id).ToRawPacket(protocol.KeepAliveOutgoingPacketId)
		s.ForEachPlayerSync(func(c *Connection) {
			c.Lock()
			if c.LastKeepAlive.ID == -1 {
				c.LastKeepAlive.ID = id
				c.LastKeepAlive.Deadline = time.Now().Add(time.Second * time.Duration(30))
				c.Write(packet)
			} else {
				if c.LastKeepAlive.Deadline.Before(time.Now()) {
					c.Disconnect("Timed out")
				}
			}
			c.Unlock()
		})
	}
}

// load loads everything the server needs.
// It is apart from the function from Start()
// because we may need to reload the server later on.
func (s *Server) load() {
	s.loadFavicon()
	s.loadBanList()
}

// loadBanList loads the players banned from the server.
func (s *Server) loadBanList() {
	s.BanList = player.NewBanList()
	if b, _ := util.Exists(banListFile); b {
		s.BanList.LoadFile(banListFile)
	}
}

// loadFavicon loads server's favicon that appears in the
// server list.
func (s *Server) loadFavicon() {
	if b, _ := util.Exists(faviconFile); b {
		contents, err := ioutil.ReadFile(faviconFile)
		if err != nil {
			log.Error("Could not load", faviconFile, err)
		} else {
			s.favicon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(contents)
		}
	}
}

// Stop stops the server.
func (s *Server) Stop() {
	s.run = false
	s.ticker.Stop()
	s.BanList.SaveFile(banListFile)
	close(s.ExitChan)
}

// handleConnection handles new connections.
func (s *Server) handleConnection(conn net.Conn) {
	c := NewConnection(conn, s)
	AssignHandler(c)
	go c.write()
	for c.IsConnected() {
		read, err := c.Next()

		if err != nil {
			// just exit
			break
		}

		if read != nil {
			c.PacketHandler.callHandler(read, c)
			read.Release()
		}
	}
	c.Disconnect("")
	// if the last connection state was the play state, we want to log his disconnection
	if c.ConnectionState == PlayState {
		s.playerLock.Lock()
		delete(s.clients, c.Player.UUID)
		s.playerLock.Unlock()

		// broadcast
		message := c.Player.Name + " has left the server."
		s.BroadcastMessage(message, protocol.DefaultMessageMode)
		log.Info(message)
	}
}

// CanConnect returns true if the given user can connect to
// the server. Otherwise, returns false and the reason why
// the player cannot connect.
func (s *Server) CanConnect(username, uuid string) (bool, string) {
	if !util.IsValidUsername(username) {
		return false, "Your username is invalid."
	}

	if banned, reason := s.BanList.IsBanned(uuid); banned {
		return false, reason
	}

	if ok, _ := s.GetPlayerByUUID(uuid); ok {
		return false, "You already logged in with this account."
	}

	return true, ""
}

// FinishLogin handles the end of player's connection to
// the server.
func (s *Server) FinishLogin(profile player.PlayerProfile, connection *Connection) {
	// TODO: Load permissions
	pl := player.Player{
		Name:        profile.Name,
		UUID:        profile.UUID,
		Permissions: make(map[string]bool),
		Profile:     profile,
		Settings:    &player.ClientSettings{},
		Location: &world.Location{
			Location3f: world.Location3f{
				X: 0,
				Y: 80,
				Z: 0,
			},
			Orientation: world.Orientation{
				Yaw:   90,
				Pitch: 0,
			},
		},
	}
	connection.Player = &pl
	s.playerLock.Lock()
	s.clients[pl.UUID] = connection
	s.playerLock.Unlock()
	packet := protocol.NewResponse()
	// send position and look packet
	packet.WriteDouble(float64(pl.Location.X)).
		WriteDouble(float64(pl.Location.Y)).
		WriteDouble(float64(pl.Location.Z))
	packet.WriteFloat(pl.Location.Yaw)
	packet.WriteFloat(pl.Location.Pitch)
	packet.WriteByte(0)
	packet.WriteVarint(int32(rand.Intn(0xFFFE)))
	connection.Write(packet.ToRawPacket(protocol.PlayerPositionAndLookPacketId))
	packet.Clear()
	// send abilities packet
	packet.WriteByte(0)
	packet.WriteFloat(1)
	packet.WriteFloat(1)
	connection.Write(packet.ToRawPacket(protocol.PlayerAbilitiesPacketId))
	packet.Clear()
	// broadcast player list item

	// announce login in chat and logs
	message := profile.Name + " has joined the server."
	log.Info(message)
	s.BroadcastMessage(message, protocol.DefaultMessageMode)
}

// GetOnlinePlayersCount returns the online players count.
func (s *Server) GetOnlinePlayersCount() uint {
	defer s.playerLock.Unlock()
	s.playerLock.Lock()
	return uint(len(s.clients))
}

// GetAllPlayers returns a slice containing all the online
// players.
func (s *Server) GetAllPlayers() []*player.Player {
	i := 0
	s.playerLock.Lock()
	ret := make([]*player.Player, len(s.clients))
	for _, client := range s.clients {
		ret[i] = client.Player
		i++
	}
	s.playerLock.Unlock()
	return ret
}

// GetPlayerByName returns true if the player associated to the given username has been found
// with the player in himself. Otherwise returns false and nil.
func (s *Server) GetPlayerByName(username string) (bool, *Connection) {
	defer s.playerLock.Unlock()
	s.playerLock.Lock()
	for _, conn := range s.clients {
		if conn.Player.Name == username {
			return true, conn
		}
	}
	return false, nil
}

// GetPlayerByUUID returns true if the player associated to the given UUID has been found
// with the player himself. Otherwise returns false and nil.
func (s *Server) GetPlayerByUUID(uuid string) (bool, *Connection) {
	defer s.playerLock.Unlock()
	s.playerLock.Lock()
	conn, ok := s.clients[uuid]
	if ok {
		return ok, conn
	} else {
		return ok, nil
	}
}

// ForEachPlayer executes the given action for each online player.
// This function runs a go routine for each player, ands waits
// the end of each routine.
func (s *Server) ForEachPlayer(action func(*Connection)) {
	s.playerLock.Lock()
	wg := sync.WaitGroup{}
	for _, client := range s.clients {
		wg.Add(1)
		go func() {
			action(client)
			wg.Done()
		}()
	}
	wg.Wait()
	s.playerLock.Unlock()
}

// ForEachPlayerSync executes the given action for each online player.
func (s *Server) ForEachPlayerSync(action func(*Connection)) {
	s.playerLock.Lock()
	for _, client := range s.clients {
		action(client)
	}
	s.playerLock.Unlock()
}

// Broadcasts the given packet to all the online players (async).
func (s *Server) BroadcastPacket(packet *protocol.RawPacket) {
	s.ForEachPlayerSync(func(c *Connection) {
		c.Write(packet)
	})
}

// Broadcasts the given message to all the players (async).
// Message's mode depends on what you want to send:
// - ChatMessageMode (mode 0): used for players only;
// - DefaultMessageMode (mode 1): what you should use (system messages);
// - ActionBarMode (mode 2): if you want to send messages above the hotbar, use this mode.
func (s *Server) BroadcastMessage(message string, mode protocol.MessageMode) {
	s.ForEachPlayerSync(func(c *Connection) {
		c.SendMessage(message, mode)
	})
}
