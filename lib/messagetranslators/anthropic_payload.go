package messagetranslators

// decodeAnthropicPayload is the source boundary for completed Anthropic
// payloads. It decodes the body once and performs optional warning inspection;
// conversion and validation of consumed core fields remain destination-specific.
func decodeAnthropicPayload(body []byte, response bool, cfg config) (map[string]any, error) {
	root, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectAnthropicObject(root, response, cfg, "$")
	return root, nil
}
