package models

import "time"

type Product struct {
	ID                        int                        `json:"id"`
	PartNumber                string                     `json:"partNumber"`
	ProductName               string                     `json:"productName"`
	IsActive                  bool                       `json:"isActive"`
	CreatedAt                 time.Time                  `json:"createdAt"`
	UpdatedAt                 time.Time                  `json:"updatedAt"`
	ProductTestConfigurations []ProductTestConfiguration `json:"productTestConfigurations"`
}
