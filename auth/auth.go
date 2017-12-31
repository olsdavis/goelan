package auth

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	. "github.com/olsdavis/goelan/player"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/olsdavis/goelan/util"
)

const (
	baseUrl = "https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%v&serverId=%v"
)

// Authenticates the current player.
func Auth(username string, sharedSecret, publicKey []byte) (*PlayerProfile, error) {
	digest := authDigest(sharedSecret, publicKey)
	resp, err := http.Get(fmt.Sprintf(baseUrl, username, digest))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var profile PlayerProfile
	err = json.Unmarshal(content, &profile)
	if err != nil {
		return nil, err
	}
	realUniqueId, err := util.StringToUUID(util.ToHypenUUID(profile.UUID))
	if err != nil {
		panic(err)
	}
	profile.RealUUID = realUniqueId
	return &profile, nil
}

// THE FOLLOWING CODE HAS BEEN TAKEN FROM: https://gist.github.com/toqueteos/5372776

// AuthDigest computes a special SHA-1 digest required for Minecraft web
// authentication on Premium servers (online-mode=true).
// Source: http://wiki.vg/Protocol_Encryption#Server
//
// Also many, many thanks to SirCmpwn and his wonderful gist (C#):
// https://gist.github.com/SirCmpwn/404223052379e82f91e6
func authDigest(sharedSecret, publicKey []byte) string {
	h := sha1.New()
	h.Write(sharedSecret)
	h.Write(publicKey)
	hash := h.Sum(nil)

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Trim away zeroes
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// little endian
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
