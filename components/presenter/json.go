package presenter

import "encoding/json"

func MarshalJSON(value any) ([]byte, error) {
	return json.Marshal(value)
}

func MarshalError(code, message string) ([]byte, error) {
	return MarshalJSON(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
