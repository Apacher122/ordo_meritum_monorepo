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
	Uid        *string
	JobID      *int
	Service    *string
	DocType    *string
	ErrorLevel *zerolog.Event
	ErrorCode  *string
	Error      error
}

// InfoLog logs a message with the service name, user ID, job ID, and doc type (if provided)
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

// ErrorLog logs a message with the service name, user ID, job ID, and doc type (if provided)
// It also logs the error code and error message if provided. If no error message is provided, it will use the default error message from error_messages.
func (l ErrorLoggerType) ErrorLog() {
	var event *zerolog.Event
	if l.ErrorLevel == log.Warn() {
		event = log.Warn()
	} else {
		event = log.Error()
	}

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
