// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return nil
	}

	var decoder Decoder
	contentType := res.Header.Get("content-type")
	if t, ok := decoderTypes[contentType]; ok {
		decoder = t(res.Body)
	} else {
		scn := bufio.NewScanner(res.Body)
		scn.Buffer(nil, bufio.MaxScanTokenSize<<9)
		decoder = &eventStreamDecoder{rc: res.Body, scn: scn}
	}
	return decoder
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder){}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for s.scn.Scan() {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
			}
			return true
		}

		// Split a string like "event: bar" into name="event" and value=" bar".
		name, value, _ := bytes.Cut(txt, []byte(":"))

		// Consume an optional space after the colon if it exists.
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			// An empty line in the for ": something" is a comment and should be ignored.
			continue
		case "event":
			event = string(value)
		case "data":
			_, s.err = data.Write(value)
			if s.err != nil {
				break
			}
			_, s.err = data.WriteRune('\n')
			if s.err != nil {
				break
			}
		}
	}

	if s.scn.Err() != nil {
		s.err = s.scn.Err()
	}

	return false
}

func (s *eventStreamDecoder) Event() Event {
	return s.evt
}

func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

func (s *eventStreamDecoder) Err() error {
	return s.err
}

type Stream[T any] struct {
	decoder Decoder
	cur     T
	err     error
}

func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
//
//		for stream.Next() {
//			data := stream.Current()
//		}
//
//	 	if stream.Err() != nil {
//			...
//	 	}
func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return false
	}

	for s.decoder.Next() {
		switch s.decoder.Event().Type {
		case "end":
			return false
		case "error":
			s.err = fmt.Errorf("received error while streaming: %s", string(s.decoder.Event().Data))
			return false
		case "data":
			var nxt T
			s.err = json.Unmarshal(s.decoder.Event().Data, &nxt)
			if s.err != nil {
				return false
			}
			s.cur = nxt
			return true
		}

		var nxt T
		s.err = json.Unmarshal(s.decoder.Event().Data, &nxt)
		if s.err != nil {
			return false
		}
		s.cur = nxt
		return true
	}

	// decoder.Next() may be false because of an error
	s.err = s.decoder.Err()

	return false
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		// already closed
		return nil
	}
	return s.decoder.Close()
}

// Patch is one JSON Patch operation of a streamed document.
type Patch struct {
	// Op is the JSON Patch operation, "add" when empty.
	Op    string
	Path  string
	Value []byte
}

// AutoPagerConfig tells an AutoPager where to find the items and the next page cursor.
type AutoPagerConfig[E any] struct {
	// Patches returns the JSON Patch operations of an event.
	Patches func(E) []Patch
	// ItemsPath is the JSON pointer to the items array, e.g. "/items".
	ItemsPath string
	// NextCursorPath is the JSON pointer to the next page cursor, e.g. "/next_cursor".
	NextCursorPath string
}

// AutoPager yields the items of a streamed query and follows the next page
// cursor across pages.
//
//	iter := client.Runs.Streaming.QueryV2StreamingAutoPaging(ctx, params)
//	for iter.Next() {
//		run := iter.Current()
//	}
//	if err := iter.Err(); err != nil {
//		...
//	}
//
// Items are yielded as they arrive while the server only appends them
// ("<items>/-" or the first document). If the server inserts an item at an
// index ("add <items>/<n>"), the final order is only known at the end of the
// page, so the rest of that page is yielded in its final order once its stream
// ends. Items yielded before the first indexed insert keep their position.
//
// Patches deeper than one item ("<items>/<n>/<field>") are ignored.
type AutoPager[E any, T any] struct {
	cfg         AutoPagerConfig[E]
	fetch       func(cursor string) *Stream[E]
	stream      *Stream[E]
	streamEnded bool
	// pending are released items not yet returned by Next.
	pending []T
	// tail are the items of the current page not released yet, in document order.
	tail []T
	// released is the number of items of the current page released so far.
	released int
	// buffered is set once the current page used an indexed patch.
	buffered bool
	next     *string
	cur      T
	err      error
	done     bool
}

func NewAutoPager[E any, T any](cfg AutoPagerConfig[E], fetch func(cursor string) *Stream[E]) *AutoPager[E, T] {
	return &AutoPager[E, T]{cfg: cfg, fetch: fetch, stream: fetch("")}
}

// Next advances to the next item. It returns false after the last page or on
// error, see Err.
func (p *AutoPager[E, T]) Next() bool {
	for {
		if p.done {
			return false
		}
		if len(p.pending) > 0 {
			p.cur, p.pending = p.pending[0], p.pending[1:]
			return true
		}
		if p.err != nil {
			return false
		}
		if !p.streamEnded {
			if p.stream.Next() {
				if p.err = p.apply(p.stream.Current()); p.err != nil {
					// Items decoded before the bad patch are still delivered, ahead of the error.
					p.release(false)
					p.stream.Close()
				}
				continue
			}
			p.err = p.stream.Err()
			p.stream.Close()
			if p.err != nil {
				return false
			}
			p.streamEnded = true
			p.release(true)
			continue
		}
		if p.next == nil {
			p.done = true
			return false
		}
		cursor := *p.next
		p.next = nil
		p.tail, p.released, p.buffered, p.streamEnded = nil, 0, false, false
		p.stream = p.fetch(cursor)
	}
}

func (p *AutoPager[E, T]) Current() T {
	return p.cur
}

func (p *AutoPager[E, T]) Err() error {
	return p.err
}

// Close releases the open stream.
func (p *AutoPager[E, T]) Close() error {
	p.done = true
	p.pending = nil
	p.tail = nil
	return p.stream.Close()
}

func (p *AutoPager[E, T]) apply(event E) error {
	for _, op := range p.cfg.Patches(event) {
		switch {
		case op.Path == "":
			var doc map[string]json.RawMessage
			if err := json.Unmarshal(op.Value, &doc); err != nil {
				return fmt.Errorf("query stream: decode document: %w", err)
			}
			if raw, ok := lookupPointer(op.Value, p.cfg.NextCursorPath); ok {
				p.next = nil
				if err := json.Unmarshal(raw, &p.next); err != nil {
					return fmt.Errorf("query stream: decode next cursor: %w", err)
				}
			} else {
				p.next = nil
			}
			if raw, ok := lookupPointer(op.Value, p.cfg.ItemsPath); ok {
				var raws []json.RawMessage
				if err := json.Unmarshal(raw, &raws); err != nil {
					return fmt.Errorf("query stream: decode items: %w", err)
				}
				items := make([]T, 0, len(raws))
				for _, raw := range raws {
					item, err := decodeItem[T](raw)
					if err != nil {
						return err
					}
					items = append(items, item)
				}
				p.tail = items[min(p.released, len(items)):]
			}
		case op.Path == p.cfg.NextCursorPath:
			p.next = nil
			if err := json.Unmarshal(op.Value, &p.next); err != nil {
				return fmt.Errorf("query stream: decode next cursor: %w", err)
			}
		default:
			rest, ok := strings.CutPrefix(op.Path, p.cfg.ItemsPath+"/")
			if !ok {
				continue
			}
			if err := p.applyItem(op, rest); err != nil {
				return err
			}
		}
	}
	if p.next != nil && *p.next == "" {
		p.next = nil
	}
	p.release(false)
	return nil
}

// applyItem applies a patch on "<items>/<rest>": "-" appends, a number is an index
// into the document. Anything else is deeper than one item and is ignored.
func (p *AutoPager[E, T]) applyItem(op Patch, rest string) error {
	isAdd := op.Op == "" || op.Op == "add"
	if rest == "-" {
		if !isAdd {
			return nil
		}
		item, err := decodeItem[T](op.Value)
		if err != nil {
			return err
		}
		p.tail = append(p.tail, item)
		return nil
	}
	index, ok := parseIndex(rest)
	if !ok {
		return nil
	}
	// Items that were already released cannot change or move: an index in front of
	// the unreleased items applies to the first of them.
	at := index - p.released
	switch {
	case isAdd:
		item, err := decodeItem[T](op.Value)
		if err != nil {
			return err
		}
		p.buffered = true
		at = min(max(at, 0), len(p.tail))
		var zero T
		p.tail = append(p.tail, zero)
		copy(p.tail[at+1:], p.tail[at:])
		p.tail[at] = item
	case op.Op == "replace" && at >= 0 && at < len(p.tail):
		item, err := decodeItem[T](op.Value)
		if err != nil {
			return err
		}
		p.buffered = true
		p.tail[at] = item
	case op.Op == "remove" && at >= 0 && at < len(p.tail):
		p.buffered = true
		p.tail = append(p.tail[:at], p.tail[at+1:]...)
	}
	return nil
}

// release moves the unreleased items to pending, unless the page is buffered and
// its stream has not ended.
func (p *AutoPager[E, T]) release(streamEnded bool) {
	if p.buffered && !streamEnded {
		return
	}
	p.pending = append(p.pending, p.tail...)
	p.released += len(p.tail)
	p.tail = nil
}

func decodeItem[T any](raw []byte) (T, error) {
	var item T
	if err := json.Unmarshal(raw, &item); err != nil {
		return item, fmt.Errorf("query stream: decode item: %w", err)
	}
	return item, nil
}

// parseIndex parses a JSON pointer array index.
func parseIndex(s string) (int, bool) {
	if s == "" || len(s) > 9 {
		return 0, false
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

// lookupPointer returns the value at the JSON pointer of a document.
func lookupPointer(doc []byte, pointer string) ([]byte, bool) {
	cur := json.RawMessage(doc)
	for _, part := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
		if part == "" {
			continue
		}
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(cur, &obj); err != nil {
			return nil, false
		}
		next, ok := obj[part]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}
