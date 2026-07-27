package messagetranslators

import "github.com/langchain-ai/langsmith-go/lib/messagetranslators/internal/utils"

func at(path string, err error) error { return utils.At(path, err) }

func decodeObject(body []byte) (map[string]any, error) { return utils.DecodeObject(body) }
func encode(v any) ([]byte, error)                     { return utils.Encode(v) }
func str(v any) (string, bool)                         { return utils.String(v) }
func arr(v any) ([]any, bool)                          { return utils.Array(v) }
func obj(v any) (map[string]any, bool)                 { return utils.Object(v) }

func requiredString(o map[string]any, key, path string) (string, error) {
	return utils.RequiredString(o, key, path)
}

func copyIf(dst, src map[string]any, to, from string) { utils.CopyIf(dst, src, to, from) }

func integer(v any, path string, positive bool) (int64, error) {
	return utils.Integer(v, path, positive)
}

func numberInRange(v any, path string, min, max float64) error {
	return utils.NumberInRange(v, path, min, max)
}

func validateOptionalBool(o map[string]any, key, path string) error {
	return utils.ValidateOptionalBool(o, key, path)
}

func parseArguments(v any, path string) (any, error) { return utils.ParseArguments(v, path) }
func parseArgumentsRaw(raw []byte, path string) (any, error) {
	return utils.ParseArgumentsRaw(raw, path)
}

func destinationID(prefix, source string, index int) string {
	return utils.DestinationID(prefix, source, index)
}
