package formatters

import (
	"errors"

	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
)

func NewEducationInfoFromPayload(payload *requests.EducationInfoPayload) (*domain.EducationInfo, error) {
	if payload == nil {
		return nil, errors.New("education payload cannot be nil")
	}
	if payload.School == "" {
		return nil, errors.New("school is a required field")
	}
	if payload.Degree == "" {
		return nil, errors.New("degree is a required field")
	}
	if payload.StartEnd == "" {
		return nil, errors.New("start_end is a required field")
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
