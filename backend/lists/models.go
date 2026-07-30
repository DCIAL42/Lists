package lists

import (
	"gorm.io/gorm"
)

const (
	ItemTypeAlbum = "album"
	ItemTypeMovie = "movie"
)

type ListItem struct {
	gorm.Model
	ListID     uint
	Type       string
	ExternalID string `json:"external_id"`
}

type List struct {
	gorm.Model
	Title     string     `json:"title"`
	CreatedBy string     `json:"created_by"`
	Items     []ListItem `gorm:"foreignKey:ListID" json:"items"`
}

type ListResponse struct {
	Title     string `json:"title"`
	CreatedBy string `json:"created_by"`
	Items     []any  `json:"items"`
}
