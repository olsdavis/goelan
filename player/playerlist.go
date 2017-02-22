package player

import (
	"encoding/json"
	"io/ioutil"
)

/*

Ban list: uuid, reason (omitempty)

Op list: uuid

*/

type BanList struct {
	players map[string]string
}

// used for JSON
type jsonBanList struct {
	profiles []jsonProfile `json:"profiles"`
}

type jsonProfile struct {
	uuid   string `json:"uuid"`
	reason string `json:"reason,omitempty"`
}

func NewBanList() *BanList {
	return &BanList{
		make(map[string]string),
	}
}

// Loads the list from the given file.
func (list *BanList) LoadFile(path string) (err error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	var bans jsonBanList
	err = json.Unmarshal(content, &bans)

	if err != nil {
		return err
	}

	for _, profile := range bans.profiles {
		list.players[profile.uuid] = profile.reason
	}

	return
}

// Returns the UUIDs that the list contains. (Can be empty.)
func (list *BanList) GetPlayers() []string {
	ret := make([]string, 0, len(list.players))
	for player, _ := range list.players {
		ret = append(ret, player)
	}
	return ret
}

// Returns true if the list has the given UUID.
func (list *BanList) HasPlayer(uuid string) bool {
	if _, ok := list.players[uuid]; ok {
		return true
	}
	return false
}
