package formatters

import (
	"errors"

	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
)

func NewEducationInfoFromPayload(payload *requests.EducationInfoPayload) (*domain.EducationInfo, *error_messages.ErrorBody) {
	if payload == nil {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_INVALID_REQUEST_FORMAT, ErrMsg: errors.New("education payload cannot be nil")}
	}
	if payload.School == "" {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_INVALID_REQUEST_FORMAT, ErrMsg: errors.New("school is a required field")}
	}
	if payload.Degree == "" {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_INVALID_REQUEST_FORMAT, ErrMsg: errors.New("degree is a required field")}
	}
	if payload.StartEnd == "" {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_INVALID_REQUEST_FORMAT, ErrMsg: errors.New("start_end is a required field")}
	}

	return &domain.EducationInfo{
		CourseWork: payload.CourseWork,
		Degree:     payload.Degree,
		Location:   payload.Location,
		School:     payload.School,
		StartEnd:   payload.StartEnd,
		GPA:        payload.GPA,
		Honors:     payload.Honors,
	}, nil
}
