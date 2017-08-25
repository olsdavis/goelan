package player

import (
	"encoding/json"
	"github.com/olsdavis/goelan/util"
	"io/ioutil"
	"os"
)

/*

Ban list: uuid, reason (omitempty)

Op list: uuid

*/

type BanList struct {
	players map[string]string
}

type BanEntry struct {
	UUID   string `json:"uuid"`
	Reason string `json:"reason,omitempty"`
}

func NewBanList() *BanList {
	return &BanList{
		make(map[string]string),
	}
}

// LoadFile loads the ban list from the given file.
func (list *BanList) LoadFile(path string) error {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	bans := make([]BanEntry, 0)
	err = json.Unmarshal(content, &bans)

	if err != nil {
		return err
	}

	for _, profile := range bans {
		list.players[profile.UUID] = profile.Reason
	}

	return nil
}

// SaveFile saves the ban list in the given file.
func (list *BanList) SaveFile(path string) error {
	content, err := json.Marshal(list)
	if err != nil {
		return err
	}
	if ok, _ := util.Exists(path); ok {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}

// AddPlayer adds the given player to the list.
func (list *BanList) AddPlayer(uuid, reason string) {
	list.players[uuid] = reason
}

// RemovePlayer removes the given player from the list.
func (list *BanList) RemovePlayer(uuid string) {
	delete(list.players, uuid)
}

// GetPlayers returns the UUIDs that the list contains. (The returned array can be empty.)
func (list *BanList) GetPlayers() []string {
	ret := make([]string, len(list.players))
	i := 0
	for _, player := range list.players {
		ret[i] = player
		i++
	}
	return ret
}

// HasPlayer returns true if the list has the given UUID.
func (list *BanList) HasPlayer(uuid string) bool {
	if _, ok := list.players[uuid]; ok {
		return true
	}
	return false
}
