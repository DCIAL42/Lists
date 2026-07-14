package client

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"strconv"
)

func NewClient(
	httpClient *http.Client,
	baseURL, resultType string,
	configParams, headers map[string]string,
	readToSearchResults func(*http.Response) (SearchResult, error),
	buildURL func(string, map[string]string) string,
	fetchToken func(context.Context, *Client),
) *Client {
	return &Client{
		httpClient,
		baseURL,
		resultType,
		configParams,
		headers,
		readToSearchResults,
		buildURL,
		fetchToken,
	}
}

func (c *Client) GetType() string {
	return c.resultType
}

func (c *Client) SetHeader(k, v string) {
	c.headers[k] = v
}

func (c *Client) GetClient() *http.Client {
	return c.httpClient
}

func nextPage(originalURL string) string {
	u, err := url.Parse(originalURL)
	if err != nil {
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

func (c *Client) tryRequest(ctx context.Context, url string) (*http.Response, error) {
	for range 3 {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			return nil, err
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(req)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 401 {
			c.fetchToken(ctx, c)
			continue
		}

		return resp, nil
	}
	return nil, errors.New("Unable to make request")
}

func (c *Client) Search(ctx context.Context, params map[string]string) (SearchResult, error) {
	maps.Copy(params, c.configParams)

	url := c.buildURL(c.baseURL, params)

	resp, err := c.tryRequest(ctx, url)

	if err != nil {
		return SearchResult{}, err
	}

	defer resp.Body.Close()

	results, err := c.readToSearchResult(resp)

	if err != nil {
		return SearchResult{}, err
	}

	var originalURL string = ctx.Value("originalURL").(string)
	results.Next = nextPage(originalURL)

	return results, nil
}
