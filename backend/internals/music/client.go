package music

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/DCIAL42/media/internals/client"
	"golang.org/x/sync/errgroup"
)

func readToSearchResult(resp *http.Response) ([]client.SearchResult, error) {
	var data Response

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return []client.SearchResult{}, err
	}

	albums := make([]client.SearchResult, 0, len(data.Results))
	g, ctx := errgroup.WithContext(resp.Request.Context())
	var mu sync.Mutex

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, r := range data.Results {
		g.Go(func() error {
			var artist string
			if len(r.Credits) > 0 {
				artist = r.Credits[0].Name
			}

			coverURL := fmt.Sprintf("https://coverartarchive.org/release-group/%s/front", r.Id)

			req, err := http.NewRequestWithContext(ctx, "GET", coverURL, nil)

			if err != nil {
				return err
			}

			resp, err := c.Do(req)

			if err != nil {
				return err
			}

			if resp.StatusCode != 200 {
				return nil
			}

			defer resp.Body.Close()

			mu.Lock()
			albums = append(albums, Album{
				Id:     r.Id,
				Title:  r.Title,
				Artist: artist,
				Cover:  coverURL,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return []client.SearchResult{}, err
	}

	return albums, nil
}

func NewMusicClient(httpClient *http.Client) *client.Client {
	return client.NewClient(
		httpClient,
		"https://musicbrainz.org/ws/2/release-group",
		"query",
		map[string]string{
			"fmt": "json",
		},
		map[string]string{
			"User-Agent": "test/1.0 ( dial1764@proton.me )",
		},
		readToSearchResult,
	)
}
