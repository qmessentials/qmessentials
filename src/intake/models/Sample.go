package models

type Sample struct {
	ID           int    `json:"id"`
	SerialNumber string `json:"serialNumber"`
	PartNumber   string `json:"partNumber"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
