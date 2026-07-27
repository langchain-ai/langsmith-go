package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// DestinationID returns a stable destination identifier for a source item.
func DestinationID(prefix, source string, index int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%s\x00%d", prefix, source, index)))
	return prefix + "_" + hex.EncodeToString(h[:12])
}
