package logger

import (
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	"github.com/rs/zerolog/log"
)

func InfoLog(
	uid string,
	service string,
	message string,
) {
	log.With().
		Str("service", service).
		Str("uid", uid).
		Logger()

	log.Info().Msg(message)
}

func ErrorLog(
	uid string,
	service string,
	errorCode string,
	message string,
) {
	log.With().
		Str("service", service).
		Str("uid", uid).
		Logger()

	em := error_messages.ErrorMessage(errorCode)
	log.Error().Err(em).Msg(message)
}
