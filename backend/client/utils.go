package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
)

func NextPage(originalURL string) string {
	u, err := url.Parse(originalURL)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	q := u.Query()
	page := q.Get("page")
	page_num, err := strconv.Atoi(page)
	if err != nil {
		page_num = 0
	}
	page_num++
	q.Set("page", strconv.Itoa(page_num))
	u.RawQuery = q.Encode()
	return u.String()
}

func Search[T cmn.ExternalItem, R db.APIResponse[T], D ResponseData[T, R]](ctx context.Context, c cmn.Client, params map[string]string) (res cmn.SearchResult, err error) {
	url := c.BuildURL(params)

	resp, err := c.TryRequest(ctx, url)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	res, err = ReadToSearchResult[T, R, D](c, c.DB(), resp, ctx.Value("userID").(string))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	var originalURL string = ctx.Value("originalURL").(string)
	res.Next = NextPage(originalURL)

	return
}

type ResponseItem[T cmn.ExternalItem] interface {
	ToDBItem() T
}

type ResponseData[T cmn.ExternalItem, U db.APIResponse[T]] interface {
	Items() []U
}

func ReadToSearchResult[T cmn.ExternalItem, R db.APIResponse[T], D ResponseData[T, R]](c cmn.Client, DB *gorm.DB, resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data D
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	items, err := db.CacheItems(DB, data.Items())
	if err != nil {
		return cmn.SearchResult{}, err
	}

	results := make([]cmn.MediaResponse, 0, len(items))
	for _, item := range items {
		results = append(results, item.ToMediaResponse())
	}

	return cmn.SearchResult{Items: results}, nil
}
