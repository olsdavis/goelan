package util

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"github.com/olsdavis/goelan/log"
	"crypto/md5"
)

type UUID struct {
	MostSig              int64 // the most significant bits of the UUID
	LeastSig             int64 // the least significant bits of the UUID
}

func NameToUUID(str string) (*UUID, error) {
	hash := md5.New()
	b := hash.Sum([]byte(str))
	b[6] &= 0x0F
	b[6] |= 0x30
	b[8] &= 0x3F
	b[8] |= 0x80

	var msb byte = 0
	var lsb byte = 0
	for i := 0; i < 8; i++ {
		msb = (msb << 8) | (b[i] & 0xFF)
	}
	for i := 8; i < 16; i++ {
		lsb = (lsb << 8) | (b[i] & 0xFF)
	}
	return &UUID{
		MostSig:              int64(msb),
		LeastSig:             int64(lsb),
	}, nil
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
		MostSig:              mostSig,
		LeastSig:             leastSig,
	}

	return uuid, nil
}

func (uuid *UUID) String() string {
	return digits(byte(uuid.MostSig) >> 32, 8) + "-" +
		digits(byte(uuid.MostSig >> 16), 4) + "-" +
		digits(byte(uuid.MostSig >> 48), 4) + "-" +
		digits(byte(uuid.LeastSig >> 48), 4) + "-" +
		digits(byte(uuid.LeastSig), 12)
}

func mustDecodeHex(str string) int64 {
	val, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		log.Error(err)
	}
	return val
}

func digits(val byte, digits uint) string {
	var hi byte = 1 << (digits * 4)
	return strconv.FormatInt(int64(hi | (val & (hi - 1))), 16)
}

// ToHyphenUUID returns the uuid with the hyphens.
func ToHyphenUUID(uuid string) string {
	// 8 - 4 - 4 - 4 - 12
	return fmt.Sprintf("%v-%v-%v-%v-%v", uuid[:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])
}
