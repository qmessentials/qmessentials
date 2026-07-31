package models

type ProductMetadata struct {
	PartNumber       string                      `json:"partNumber"`
	MetadataType     *string                     `json:"metadataType"`
	MetadataSelector *string                     `json:"metadataSelector"`
	Definitions      []ProductMetadataDefinition `json:"definitions"`
}

type ProductMetadataDefinition struct {
	MetadataKey string `json:"metadataKey"`
	ValueType   string `json:"valueType"`
}
