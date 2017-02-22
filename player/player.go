package player

type Player struct {
	Name string
	UUID string
	// key => the permission; value => true if the player has the permission
	Permissions map[string]bool
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
