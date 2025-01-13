package dtos

type Photo struct {
	ID        string `json:"id"`
	Data      string `json:"data"` // Stored in b64
	IsDeleted bool   `json:"isDeleted"`
}

type PhotoDB struct {
	ID         string `json:"id"`
	DataOrigin string `json:"dataOrigin"` // Stored in b64
	Data75     string `json:"data75"`     // Stored in b64
	Data50     string `json:"data50"`     // Stored in b64
	Data25     string `json:"data25"`     // Stored in b64
	IsDeleted  bool   `json:"isDeleted"`
}
