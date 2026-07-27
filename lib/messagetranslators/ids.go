package messagetranslators

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// destinationID returns a stable destination identifier for a source item.
func destinationID(prefix, source string, index int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%s\x00%d", prefix, source, index)))
	return prefix + "_" + hex.EncodeToString(h[:12])
}
