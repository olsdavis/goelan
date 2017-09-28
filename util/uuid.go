package util

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"github.com/olsdavis/goelan/log"
)

type UUID struct {
	StringRepresentation string
	MostSig              int64 // the most significant bits of the UUID
	LeastSig             int64 // the least significant bits of the UUID
}

// StringToUUID converts the given string to a UUID struct.
func StringToUUID(str string) (*UUID, error) {
	components := strings.Split(str, "-")
	if len(components) != 5 {
		return nil, errors.New("invalid UUID string " + str)
	}

	mostSig := mustDecodeHex(components[0])
	mostSig <<= 16
	mostSig |= mustDecodeHex(components[1])
	mostSig <<= 16
	mostSig |= mustDecodeHex(components[2])

	leastSig := mustDecodeHex(components[3])
	leastSig <<= 48
	leastSig |= mustDecodeHex(components[4])

	uuid := &UUID{
		StringRepresentation: str,
		MostSig:              mostSig,
		LeastSig:             leastSig,
	}

	return uuid, nil
}

func mustDecodeHex(str string) int64 {
	val, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		log.Error(err)
	}
	return val
}

// ToHypenUUID returns the uuid with the hyphens.
func ToHypenUUID(uuid string) string {
	// 8 - 4 - 4 - 4 - 12
	return fmt.Sprintf("%v-%v-%v-%v-%v", uuid[:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])
}
