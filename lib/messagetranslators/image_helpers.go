package messagetranslators

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// anthropicImageURL validates an Anthropic image source and returns the URL
// representation shared by the Responses and Chat Completions converters.
func anthropicImageURL(o map[string]any, path string, incompleteSource error) (string, error) {
	s, ok := obj(o["source"])
	if !ok {
		return "", at(path+".source", ErrInvalidWireData)
	}
	switch s["type"] {
	case "url":
		u, ok := str(s["url"])
		if !ok || u == "" {
			return "", at(path+".source.url", ErrInvalidWireData)
		}
		return u, nil
	case "base64":
		media, mediaOK := str(s["media_type"])
		data, dataOK := str(s["data"])
		if !mediaOK || !dataOK || media == "" || data == "" {
			return "", at(path+".source", incompleteSource)
		}
		if _, err := base64.StdEncoding.DecodeString(data); err != nil {
			return "", at(path+".source.data", fmt.Errorf("%w: invalid base64", ErrInvalidWireData))
		}
		return "data:" + media + ";base64," + data, nil
	default:
		return "", at(path+".source.type", ErrUnsupported)
	}
}

// imageURLToAnthropic turns either an ordinary URL or a base64 data URL into
// an Anthropic image block. urlPath is kept explicit so callers retain their
// API-specific error paths.
func imageURLToAnthropic(u, urlPath string) (map[string]any, error) {
	if strings.HasPrefix(u, "data:") {
		parts := strings.SplitN(strings.TrimPrefix(u, "data:"), ",", 2)
		if len(parts) != 2 || !strings.HasSuffix(parts[0], ";base64") {
			return nil, at(urlPath, ErrUnsupported)
		}
		media := strings.TrimSuffix(parts[0], ";base64")
		if media == "" || parts[1] == "" {
			return nil, at(urlPath, ErrInvalidWireData)
		}
		if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
			return nil, at(urlPath, ErrInvalidWireData)
		}
		return map[string]any{"type": "image", "source": map[string]any{"type": "base64", "media_type": media, "data": parts[1]}}, nil
	}
	return map[string]any{"type": "image", "source": map[string]any{"type": "url", "url": u}}, nil
}
