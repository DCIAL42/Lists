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

func Search(ctx context.Context, c cmn.Client, params map[string]string) (res cmn.SearchResult, err error) {
	url := c.BuildURL(params)

	resp, err := c.TryRequest(ctx, url)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	res, err = c.ReadToSearchResult(resp, ctx.Value("userID").(string))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	var originalURL string = ctx.Value("originalURL").(string)
	res.Next = NextPage(originalURL)

	return
}

type ResponseItem[T db.ExternalItem] interface {
	ToDBItem() T
}

type ResponseData[T db.ExternalItem, U ResponseItem[T]] interface {
	Items() []U
}

func TestRead[T db.ExternalItem, R ResponseItem[T], D ResponseData[T, R]](DB *gorm.DB, resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data D
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	responseItems := data.Items()

	items := make([]T, 0, len(responseItems))
	mediaIDs := make([]uint, 0, len(responseItems))

	for _, r := range responseItems {
		item := r.ToDBItem()
		if _, err = db.TrySaveItem(DB, item); err != nil {
			return
		}
		items = append(items, item)
		mediaIDs = append(mediaIDs, item.GetMediaID())
	}

	var tracking []cmn.TrackingItem

	if err = DB.Where("user_id = ?", userID).Where("media_id IN ?", mediaIDs).Find(&tracking).Error; err != nil {
		return
	}

	trackingByMediaID := make(map[uint]*cmn.TrackingItem, len(tracking))
	for i := range tracking {
		trackingByMediaID[tracking[i].MediaID] = &tracking[i]
	}

	results := make([]cmn.MediaResponse, 0, len(items))

	for i := range items {
		items[i].GetMedia().Tracking = trackingByMediaID[items[i].GetMediaID()]
		results = append(results, items[i].ToMediaResponse())
	}

	return cmn.SearchResult{Items: results}, nil
}
