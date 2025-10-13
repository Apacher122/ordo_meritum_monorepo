package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			return errors.New(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Invalid type for field '%s'. Expected '%s' but received a '%s'.", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
			return errors.New(msg)

		default:
			return errors.New("Invalid request body: " + err.Error())
		}
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
