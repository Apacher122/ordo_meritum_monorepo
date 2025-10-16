package mocks

import (
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/events"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/shared/utils/formatters"
)

func GetMockDocumentEvent(uid string, jobId int, docType string) events.DocumentEvent {
	return events.DocumentEvent{
		JobID:       jobId,
		UserId:      uid,
		CompanyName: "Tech Solutions Inc.",
		DocType:     docType,
		UserInfo: requests.UserInfoPayload{
			FirstName:       "Jane",
			LastName:        "Doe",
			CurrentLocation: "San Francisco, CA",
			Email:           "jane.doe@email.com",
			Github:          "github.com/janedoe",
			Linkedin:        "linkedin.com/in/janedoe",
			Mobile:          "555-123-4567",
			Summary:         "A highly motivated software engineer with 5 years of experience in building scalable web applications.",
		},
		EducationInfo: requests.EducationInfoPayload{
			CourseWork: formatters.StringToPtr("Advanced Algorithms, Database Systems, Web Development"),
			Degree:     "Bachelor of Science in Computer Science",
			Location:   "Berkeley, CA",
			School:     "University of California, Berkeley",
			StartEnd:   "2015-2019",
		},
		Resume: domain.Resume{
			Summary: []domain.SummaryBody{
				{
					Sentence:               "Skilled in full-stack development with a focus on React and Go.",
					JustificationForChange: "Added more specific keywords.",
					NewSuggestion:          true,
				},
				{
					Sentence:               "Proven ability to lead projects from conception to completion.",
					JustificationForChange: "BLahBlah", NewSuggestion: true},
			},
			Skills: []domain.Skills{
				{
					Category:                "Programming Languages",
					SkillItem:               []string{"Go", "TypeScript", "Python"},
					JustificationForChanges: "Consolidated skill categories.",
				},
				{
					Category:                "Frameworks & Libraries",
					SkillItem:               []string{"React", "Node.js", "Gin", "gorilla/mux"},
					JustificationForChanges: "adsafasdfasdf",
				},
			},
			Experiences: []domain.Experience{
				{
					BulletPoints: []domain.BulletPoint{
						{
							Text:                   "Developed and maintained microservices using Go, improving API response times by 30%.",
							IsNewSuggestion:        false,
							JustificationForChange: "blahblah",
						},
						{
							Text:                   "Engineered a new real-time notification system using WebSockets, increasing user engagement.",
							IsNewSuggestion:        true,
							JustificationForChange: "Quantified the impact of the achievement.",
						},
					},
					Company:  "Innovate Corp",
					ID:       "exp-1",
					Position: "Senior Software Engineer",
					Start:    "2021",
					End:      "Present",
				},
			},
			Projects: []domain.Project{
				{
					BulletPoints: []domain.BulletPoint{
						{
							Text:                   "Built a full-stack e-commerce platform with a React frontend and Go backend.",
							IsNewSuggestion:        false,
							JustificationForChange: "asdfasdasdf",
						},
					},
					Role:   "Lead Developer",
					ID:     "proj-1",
					Name:   "E-Commerce Platform",
					Status: "Completed",
				},
			},
		},
		CoverLetter: domain.CoverLetter{
			CompanyProperName: "Tech Solutions Incorporated",
			JobTitle:          "Software Engineer",
			Body: domain.CoverLetterBody{
				About:           "I am writing to express my interest in the Software Engineer position advertised on your website.",
				Experience:      "In my previous role at Innovate Corp, I was responsible for key microservice development.",
				WhatIBring:      "I bring a strong proficiency in Go and a passion for clean, efficient code.",
				RevisionSummary: "Revised to better align with company values.",
			},
		},
	}
}
