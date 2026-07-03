package client

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/url"
)

func NewClient(httpClient *http.Client, baseURL, resultType string, configParams, headers map[string]string, readToSearchResults func(*http.Response) ([]SearchResult, error)) *Client {
	return &Client{
		httpClient,
		baseURL,
		resultType,
		configParams,
		headers,
		readToSearchResults,
	}
}

func buildUrl(params map[string]string) string {
	query := url.Values{}

	for k, v := range params {
		query.Set(k, v)
	}

	return query.Encode()
}

func (c *Client) GetType() string {
	return c.resultType
}

func (c *Client) Search(ctx context.Context, params map[string]string) ([]SearchResult, error) {
	maps.Copy(params, c.configParams)

	q := buildUrl(params)

	fmt.Println(q)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"?"+q, nil)

	if err != nil {
		return []SearchResult{}, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return []SearchResult{}, err
	}

	defer resp.Body.Close()

	results, err := c.readToSearchResult(resp)

	if err != nil {
		return []SearchResult{}, err
	}

	return results, nil
}
