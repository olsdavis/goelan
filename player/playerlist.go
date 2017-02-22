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

type JsonProfile struct {
	UUID   string `json:"uuid"`
	Reason string `json:"reason"`
}

func NewBanList() *BanList {
	return &BanList{
		make(map[string]string),
	}
}

// Loads the list from the given file.
func (list *BanList) LoadFile(path string) error {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	bans := make([]JsonProfile, 0)
	err = json.Unmarshal(content, &bans)

	if err != nil {
		return err
	}

	for _, profile := range bans {
		list.players[profile.UUID] = profile.Reason
	}

	return nil
}

// Saves the list to the given file.
func (list *BanList) SaveFile(path string) error {
	return nil
}

// Adds the given player to the list.
func (list *BanList) AddPlayer(uuid, reason string) {
	list.players[uuid] = reason
}

// Removes the given player from the list.
func (list *BanList) RemovePlayer(uuid, reason string) {
	delete(list.players, uuid)
}

// Returns the UUIDs that the list contains. (Can be empty.)
func (list *BanList) GetPlayers() []string {
	ret := make([]string, len(list.players))
	i := 0
	for _, player := range list.players {
		ret[i] = player
		i++
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
