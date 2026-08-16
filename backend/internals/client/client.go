package client

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/DCIAL42/lists/cmn"
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

	res, err = c.ReadToSearchResult(resp)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	var originalURL string = ctx.Value("originalURL").(string)
	res.Next = NextPage(originalURL)

	return
}
