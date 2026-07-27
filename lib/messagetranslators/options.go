package messagetranslators

import "time"

// Option configures one conversion or stream converter. The zero set of options
// silently tolerates unknown source fields, uses time.Now for generated
// timestamps, and reports no usage.
type Option func(*config)

// WithWarningHandler installs h to receive schema-drift and lossy-mapping
// warnings synchronously, in source traversal order. Installing a handler is
// what enables source inspection at all; see README.md for the boundaries of
// what is inspected. h must not panic, and panics are deliberately not
// recovered.
func WithWarningHandler(h WarningHandler) Option {
	return func(c *config) { c.warnings = h }
}

// WithClock replaces time.Now as the source of generated timestamps, which the
// converters need because Anthropic payloads carry no creation time while
// Responses and Chat Completions payloads require one. Supply a clock to make
// output deterministic under test, or to align generated timestamps with the
// gateway's own request clock rather than with conversion time.
func WithClock(now func() time.Time) Option {
	return func(c *config) {
		if now != nil {
			c.now = now
		}
	}
}

// WithUsageHandler installs h to receive final token accounting when a
// conversion or stream reaches a terminal state.
func WithUsageHandler(h UsageHandler) Option {
	return func(c *config) { c.usage = h }
}

type config struct {
	warnings WarningHandler
	usage    UsageHandler
	now      func() time.Time
}

func newConfig(options []Option) config {
	c := config{now: time.Now}
	for _, option := range options {
		if option != nil {
			option(&c)
		}
	}
	return c
}

func (c config) warn(w Warning) {
	if c.warnings != nil {
		c.warnings(w)
	}
}

func (c config) unknownField(path, field string) {
	c.warn(Warning{
		Code:    WarningUnknownField,
		Path:    path,
		Field:   field,
		Message: "unknown source field " + quote(field) + " at " + path + " was ignored",
	})
}

// lossy reports source data that was understood and validated but cannot be
// represented at the destination. Unlike an unknown field, the loss is known
// and intentional; the warning exists so a gateway can measure it.
func (c config) lossy(path, field, message string) {
	c.warn(Warning{Code: WarningLossyConversion, Path: path, Field: field, Message: message})
}

func (c config) reportUsage(u Usage) {
	if c.usage != nil {
		c.usage(u)
	}
}

func quote(s string) string { return `"` + s + `"` }
