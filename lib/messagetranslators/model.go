package messagetranslators

import "fmt"

// resolveModel picks the destination model. An override wins; otherwise the
// source model carries through. The wire value is validated whenever it is
// present, so an override never hides a malformed source body.
func resolveModel(o map[string]any, override, path string) (string, error) {
	if _, present := o["model"]; present {
		wire, err := requiredString(o, "model", path)
		if err != nil {
			return "", err
		}
		if override == "" {
			return wire, nil
		}
	} else if override == "" {
		return "", at(path+".model", fmt.Errorf("%w: non-empty string required", ErrInvalidWireData))
	}
	return override, nil
}

// resolveRequiredModel is resolveModel for completed responses, where the source
// must always name the model that produced it even when the caller overrides the
// value reported downstream.
func resolveRequiredModel(o map[string]any, override, path string) (string, error) {
	wire, err := requiredString(o, "model", path)
	if err != nil {
		return "", err
	}
	if override == "" {
		return wire, nil
	}
	return override, nil
}
