package models

import (
	"time"

	"github.com/google/uuid"
)

type TestResult struct {
	ID                  uuid.UUID `json:"id"`
	SampleID            int       `json:"sampleId"`
	PartNumber          string    `json:"partNumber"`
	ProductTestSequence int       `json:"productTestSequence"`
	Modifiers           []string  `json:"modifiers"`
	TestResult          float64   `json:"testResult"`
	Unit                string    `json:"unit"`
	DecimalPlaces       int       `json:"decimalPlaces"`
	MinValue            float64   `json:"minValue"`
	MaxValue            float64   `json:"maxValue"`
	HashValue           string    `json:"hashValue"`
	VoidedAt            time.Time `json:"voidedAt"`
	VoidedBy            string    `json:"voidedBy"`
	VoidedReason        string    `json:"voidedReason"`
	VoidComment         string    `json:"voidComment"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
