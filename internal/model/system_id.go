package model

import (
	"crypto/sha256"
	"github.com/jxskiss/base62"
	"github.com/lzap/deagon"
	"hash/fnv"
	"sort"
	"strings"
)

type SystemID struct {
	// Base62-enconded 256-bit string id that is cryptographically safe, very likely unique
	Long string

	// Base62-enconded 128-bit string id that not cryptographically safe, likely unique
	Short string

	// Friendly name based on 25-bit, hopefully unique, only to be used as a "display" name
	FriendlyName string
}

var formatter = deagon.NewCapitalizedSpaceFormatter()

func NewSystemID(serial string, serialAndMACs ...string) SystemID {
	sort.Strings(serialAndMACs)
	sha := sha256.New()
	fnv32 := fnv.New32()
	fnv128 := fnv.New128()

    _, _ = sha.Write([]byte(strings.TrimSpace(serial)))
    _, _ = fnv32.Write([]byte(strings.TrimSpace(serial)))
    _, _ = fnv128.Write([]byte(strings.TrimSpace(serial)))

	for _, str := range serialAndMACs {
        buf := []byte(strings.ToLower(strings.TrimSpace(str)))
		_, _ = sha.Write(buf)
		_, _ = fnv32.Write(buf)
		_, _ = fnv128.Write(buf)
	}
	shaHash := sha.Sum(nil)
	fnvHash := fnv128.Sum(nil)

	return SystemID{
		Long:         base62.EncodeToString(shaHash),
		Short:        base62.EncodeToString(fnvHash),
		FriendlyName: deagon.Name(formatter, int(fnv32.Sum32())),
	}
}
