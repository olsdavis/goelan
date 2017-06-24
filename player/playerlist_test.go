package player

import (
	"testing"
)

var (
	banList *BanList = NewBanList()
)

func TestLoadFile(t *testing.T) {
	err := banList.LoadFile("banlist.json.test")
	if err != nil {
		t.Error("Could not retrieve the ban list from the file:", err)
	}
}

func TestGetPlayers(t *testing.T) {
	banList.LoadFile("banlist.json.test")
	players := banList.GetPlayers()
	if len(players) == 0 {
		t.Error("GetPlayers() returns an empty slice. Currently returning:", players)
	}
}
