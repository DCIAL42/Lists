package lists

import "gorm.io/gorm"

type Settings struct {
	Page  uint
	Limit uint
	Full  bool
}

type ListQueryCfg struct {
	Order string
	By    string
}

func orderByLikes(order string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Select("lists.*, COUNT(likes.id) AS like_count").
			Joins("LEFT JOIN likes ON likes.list_id = lists.id AND likes.deleted_at IS NULL").
			Group("lists.id").
			Order("like_count " + order)
	}
}

func standardOrder(by, order string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order(by + " " + order)
	}
}
