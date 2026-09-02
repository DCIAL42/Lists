package lists

import (
	"github.com/DCIAL42/lists/cmn"
)

type ListItemRequest struct {
	Type    cmn.MediaType `json:"type"`
	MediaID uint          `json:"media_id"`
}

type ListRequest struct {
	Title string            `json:"title"`
	Items []ListItemRequest `json:"items"`
}

type ListItem struct {
	cmn.Model
	ListID  uint
	MediaID uint

	Media cmn.Media `gorm:"foreignKey:MediaID"`
}

type List struct {
	cmn.Model
	UserID string
	Title  string
	Items  []ListItem `gorm:"foreignKey:ListID"`
	Cover  string
	Likes  []cmn.Like `gorm:"foreignKey:ListID"`
}

type ListMeta struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title"`
	CreatedBy string `json:"created_by"`
	Cover     string `json:"cover"`
	Likes     uint   `json:"likes"`
}

type ListResponse struct {
	ListMeta
	Items []cmn.MediaResponse `json:"items"`
}

type ListsPreviewResponse struct {
	Lists []ListMeta `json:"lists"`
	Next  string     `json:"next,omitempty"`
	Page  uint       `json:"page"`
	Count uint       `json:"count"`
}

type ListsFullResponse struct {
	Lists []ListResponse `json:"lists"`
	Next  string         `json:"next,omitempty"`
	Page  uint           `json:"page"`
	Count uint           `json:"count"`
}
