package lists

import (
	"github.com/DCIAL42/lists/cmn"
)

type ListItem struct {
	cmn.Model
	ListID     uint
	Type       cmn.MediaType `json:"type"`
	ExternalID string        `json:"external_id"`
}

type UpdateListRequest struct {
	Title string     `json:"title"`
	Items []ListItem `json:"items"`
}

type List struct {
	cmn.Model
	UserID string
	Title  string     `json:"title"`
	Items  []ListItem `gorm:"foreignKey:ListID" json:"items"`
	Cover  string     `json:"cover"`
}

type ListMeta struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title"`
	CreatedBy string `json:"created_by"`
	Cover     string `json:"cover"`
}

type ListResponse struct {
	ListMeta
	Items []cmn.MediaItem `json:"items"`
}

type ListsPreviewResponse struct {
	Lists []ListMeta `json:"lists"`
	Next  string     `json:"next,omitempty"`
	Page  uint       `json:"page"`
	Count uint       `json:"count"`
}

type ListsResponse struct {
	Lists []ListResponse `json:"lists"`
	Next  string         `json:"next,omitempty"`
	Page  uint           `json:"page"`
	Count uint           `json:"count"`
}
