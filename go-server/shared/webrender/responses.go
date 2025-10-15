package webrender

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/shared/middleware"
	error_response "github.com/ordo_meritum/shared/types/errors"
)

func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request, v T) error {
	err := json.NewDecoder(r.Body).Decode(v)

	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")

		errorPayload := error_response.ErrorResponse[struct{}]{
			ErrorCode: error_response.BAD_REQUEST,
			Message:   "The request body is malformed or unreadable.",
		}
		middleware.JSON(w, http.StatusBadRequest, errorPayload)
		return err
	}

	return nil
}
