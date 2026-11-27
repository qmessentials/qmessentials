package models

type Product struct {
	ID                        int                        `json:"id"`
	PartNumber                string                     `json:"partNumber"`
	ProductName               string                     `json:"productName"`
	IsActive                  bool                       `json:"isActive"`
	ProductTestConfigurations []ProductTestConfiguration `json:"productTestConfigurations"`
}
