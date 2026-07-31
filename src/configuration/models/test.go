package models

import "time"

type Test struct {
	ID                      int       `json:"id"`
	TestName                string    `json:"testName"`
	CanonicalTestName       string    `json:"canonicalTestName"`
	TestUnitCategory        string    `json:"testUnitCategory"`
	DocumentationReferences []string  `json:"documentationReferences"`
	AreModifiersAllowed     bool      `json:"areModifiersAllowed"`
	IsActive                bool      `json:"isActive"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}
