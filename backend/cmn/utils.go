package cmn

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(models ...any) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return
	}

	err = db.AutoMigrate(models...)

	return
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
}

func (e *HttpError) Error() string {
	return e.Message
}

func HandleError(c *gin.Context, err error) {
	if httpErr, ok := errors.AsType[*HttpError](err); ok {
		c.IndentedJSON(httpErr.Code, gin.H{"error": httpErr.Error()})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
