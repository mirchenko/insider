package model

type ListResponseMetadata struct {
	Count  int64 `json:"count"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
