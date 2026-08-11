package lists

import (
	"github.com/DCIAL42/media/cmn"
)

type ListItem struct {
	cmn.Model
	ListID     uint
	Type       cmn.MediaType `json:"type"`
	ExternalID string        `json:"external_id"`
}

type List struct {
	cmn.Model
	UserID    string
	Title     string     `json:"title"`
	CreatedBy string     `json:"created_by"`
	Items     []ListItem `gorm:"foreignKey:ListID" json:"items"`
}

type ListResponse struct {
	ID        uint            `json:"id,omitempty"`
	Title     string          `json:"title"`
	CreatedBy string          `json:"created_by"`
	Items     []cmn.MediaItem `json:"items"`
}
