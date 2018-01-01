package player

import (
	"github.com/olsdavis/goelan/world"
	"github.com/olsdavis/goelan/util"
)

type PlayerProfile struct {
	RealUUID   *util.UUID `json:"-"`
	UUID       string     `json:"id"`
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type Player struct {
	// key => the permission; value => true if the player has the permission
	Permissions map[string]bool
	Profile     PlayerProfile
	Settings    *ClientSettings
	Location    *world.Location
	GameMode    GameMode
}

// HasPermission returns true if the player has the given permission.
func (player *Player) HasPermission(permission string) bool {
	can, ok := player.Permissions[permission]
	return can && ok
}

// SetPermission sets whether the player will have the given permission or not.
func (player *Player) SetPermission(permission string, can bool) {
	player.Permissions[permission] = can
}

// GetName returns current player's name.
func (player *Player) GetName() string {
	return player.Profile.Name
}
