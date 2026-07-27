package messagetranslators

import "github.com/langchain-ai/langsmith-go/lib/messagetranslators/internal/utils"

var (
	// ErrInvalidWireData means that JSON or an event payload is malformed.
	ErrInvalidWireData = utils.ErrInvalidWireData
	// ErrUnsupported means that a known, semantically significant feature has no safe mapping.
	ErrUnsupported = utils.ErrUnsupported
	// ErrInvalidSequence means that otherwise valid stream events arrived out of order.
	ErrInvalidSequence = utils.ErrInvalidSequence
	// ErrTruncatedStream means Finish was called before a terminal stream event.
	ErrTruncatedStream = utils.ErrTruncatedStream
)

// ConversionError adds a JSON path or stream-event location to a conversion error.
type ConversionError = utils.ConversionError
