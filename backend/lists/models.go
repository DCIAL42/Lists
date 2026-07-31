package lists

import (
	"time"

	"github.com/DCIAL42/media/internals/client"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListItem struct {
	Model
	ListID     uint
	Type       client.MediaType
	ExternalID string `json:"external_id"`
}

type List struct {
	Model
	Title     string     `json:"title"`
	CreatedBy string     `json:"created_by"`
	Items     []ListItem `gorm:"foreignKey:ListID" json:"items"`
}

// type ListItemResponse struct {
// 	client.MediaItem
// }

type ListResponse struct {
	Title     string             `json:"title"`
	CreatedBy string             `json:"created_by"`
	Items     []client.MediaItem `json:"items"`
}
