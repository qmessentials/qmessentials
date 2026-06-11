package models

import "time"

type ProductTestConfiguration struct {
	ID                int       `json:"id"`
	ProductID         int       `json:"productId"`
	TestID            int       `json:"testId"`
	SpecificModifiers []string  `json:"specificModifiers"`
	Unit              string    `json:"unit"`
	DecimalPlaces     int       `json:"decimalPlaces"`
	MinValue          *float64  `json:"minValue"`
	MaxValue          *float64  `json:"maxValue"`
	IsCritical        bool      `json:"isCritical"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Test              *Test     `json:"test"`
}
