package logger

import (
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type InfoLoggerType struct {
	Uid     *string
	JobID   *int
	Service *string
	DocType *string
	Message string
}

type ErrorLoggerType struct {
	Uid       *string
	JobID     *int
	Service   *string
	DocType   *string
	ErrorCode *string
	Error     error
}

func (l InfoLoggerType) InfoLog() {
	lg := log.Info()
	if l.Service != nil {
		lg.Str("service", *l.Service)
	}

	if l.Uid != nil {
		lg.Str("uid", *l.Uid)
	}

	if l.JobID != nil {
		lg.Int("jobID", *l.JobID)
	}

	if l.DocType != nil {
		lg.Str("docType", *l.DocType)
	}

	lg.Msg(l.Message)
}

func (l ErrorLoggerType) ErrorLog() {
	event := log.Error().
		Str("uid", *l.Uid)

	if l.Uid != nil {
		event.Str("uid", *l.Uid)
	}

	if l.JobID != nil {
		event.Int("jobID", *l.JobID)
	}

	if l.Service != nil {
		event.Str("service", *l.Service)
	}

	if l.DocType != nil {
		event.Str("docType", *l.DocType)
	}

	if l.ErrorCode != nil {
		event.Str("error_code", *l.ErrorCode)
	}

	if l.Error != nil {
		event.Err(l.Error)
	} else {
		event.Err(error_messages.ErrorMessage(*l.ErrorCode))
	}
}

func ErrorLog(errorCode string, err error, event *zerolog.Event) *zerolog.Event {
	if err == nil {
		err = error_messages.ErrorMessage((errorCode))
	}
	ctx := event.Str("error_code", errorCode).Err(err)
	return ctx
}
