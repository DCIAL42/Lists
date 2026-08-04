package cmn

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dst ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(dst...)

	return db, err
}

func NextPage(u *url.URL) string {
	q := u.Query()
	page := q.Get("page")
	page_num, err := strconv.Atoi(page)
	if err != nil {
		page_num = 1
	}
	page_num++
	q.Set("page", strconv.Itoa(page_num))
	u.RawQuery = q.Encode()
	return u.String()
}

// func ParseID(idStr string) (uint, error) {
// 	val, err := strconv.ParseUint(idStr, 10, 64)
//
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	id := uint(val)
//
// 	return id, nil
// }

func ParseParam[T int | uint](c *gin.Context, p string) (val T, err error) {
	v := c.Param(p)

	switch any(val).(type) {
	case int:
		x, err := strconv.Atoi(v)
		return T(x), err
	case uint:
		x, err := strconv.ParseUint(v, 10, 64)
		return T(x), err
	}

	err = errors.New("unsupported type")

	return
}

type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}
