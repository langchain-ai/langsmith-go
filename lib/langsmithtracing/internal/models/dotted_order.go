package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NewDottedSegment creates the first segment of a dotted order string.
// Format: YYYYMMDDTHHMMSSffffffZ<uuid> (no dot before microseconds).
func NewDottedSegment(t time.Time, id uuid.UUID) string {
	utc := t.UTC()
	return fmt.Sprintf("%s%06dZ%s", utc.Format("20060102T150405"), utc.Nanosecond()/1000, id.String())
}

// AppendDotted appends a new segment to an existing dotted order.
func AppendDotted(parent string, t time.Time, id uuid.UUID) string {
	return parent + "." + NewDottedSegment(t, id)
}
