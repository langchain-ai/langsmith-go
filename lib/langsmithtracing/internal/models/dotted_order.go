package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NewDottedSegment creates the first segment of a dotted order string.
func NewDottedSegment(t time.Time, id uuid.UUID) string {
	return fmt.Sprintf("%s%s", t.UTC().Format("20060102T150405.000000"), "Z"+id.String())
}

// AppendDotted appends a new segment to an existing dotted order.
func AppendDotted(parent string, t time.Time, id uuid.UUID) string {
	return parent + "." + NewDottedSegment(t, id)
}
