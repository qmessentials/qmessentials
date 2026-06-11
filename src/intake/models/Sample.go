package models

import "time"

type Sample struct {
	ID           int       `json:"id"`
	SerialNumber string    `json:"serialNumber"`
	PartNumber   string    `json:"partNumber"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
