package player

import "github.com/olsdavis/goelan/world"

type PlayerProfile struct {
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
	Name string
	UUID string
	// key => the permission; value => true if the player has the permission
	Permissions map[string]bool
	Profile     PlayerProfile
	Settings    *ClientSettings
	Location    *world.Location
}

// Returns true if the player has the given permission.
func (player *Player) HasPermission(permission string) bool {
	can, ok := player.Permissions[permission]
	return can && ok
}

// If can is true, the player will have the given permission; otherwise, he won't.
func (player *Player) SetPermission(permission string, can bool) {
	player.Permissions[permission] = can
}
