package messagetranslators

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidWireData means that JSON or an event payload is malformed.
	ErrInvalidWireData = errors.New("messagetranslators: invalid wire data")
	// ErrUnsupported means that a known, semantically significant feature has no safe mapping.
	ErrUnsupported = errors.New("messagetranslators: unsupported feature")
	// ErrInvalidSequence means that otherwise valid stream events arrived out of order.
	ErrInvalidSequence = errors.New("messagetranslators: invalid stream sequence")
	// ErrTruncatedStream means Finish was called before a terminal stream event.
	ErrTruncatedStream = errors.New("messagetranslators: truncated stream")
)

// ConversionError adds a JSON path or stream-event location to a conversion error.
type ConversionError struct {
	Path string
	Err  error
}

func (e *ConversionError) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("messagetranslators: %v", e.Err)
	}
	return fmt.Sprintf("messagetranslators: %s: %v", e.Path, e.Err)
}

func (e *ConversionError) Unwrap() error { return e.Err }

func at(path string, err error) error { return &ConversionError{Path: path, Err: err} }
