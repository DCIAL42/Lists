package music

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/DCIAL42/media/internals/client"
)

type PopPayload struct {
	Mbids []string `json:"release_group_mbids"`
}

type PopResponse struct {
	Id          string `json:"release_group_mbid"`
	ListenCount int    `json:"total_listen_count"`
	UserCount   int    `json:"total_user_count"`
}

func fetchPops(albums []Album, ctx context.Context) ([]PopResponse, error) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	ids := make([]string, 0, len(albums))

	for _, a := range albums {
		ids = append(ids, a.Id)
	}

	idData := PopPayload{Mbids: ids}
	jsonData, err := json.Marshal(idData)
	if err != nil {
		return []PopResponse{}, err
	}
	body := bytes.NewBuffer(jsonData)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.listenbrainz.org/1/popularity/release-group", body)

	if err != nil {
		return []PopResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", "a344070c-0efb-47df-8b0a-52347fc19fae"))

	resp, err := c.Do(req)

	if err != nil {
		return []PopResponse{}, err
	}

	defer resp.Body.Close()

	var pops []PopResponse

	err = json.NewDecoder(resp.Body).Decode(&pops)

	return pops, nil
}

func readToSearchResult(resp *http.Response) ([]client.SearchResult, error) {
	var data Response

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return []client.SearchResult{}, err
	}

	albums := make([]Album, 0, len(data.Results))

	for _, r := range data.Results {
		var artist string
		if len(r.Credits) > 0 {
			artist = r.Credits[0].Name
		}

		coverURL := fmt.Sprintf("https://coverartarchive.org/release-group/%s/front", r.Id)

		albums = append(albums, Album{
			Id:     r.Id,
			Title:  r.Title,
			Artist: artist,
			Cover:  coverURL,
		})
	}

	pops, err := fetchPops(albums, resp.Request.Context())

	if err != nil {
		return []client.SearchResult{}, err
	}

	for i := range albums {
		albums[i].ListenCount = pops[i].ListenCount
	}

	sort.Slice(albums, func(i, j int) bool {
		return albums[i].ListenCount > albums[j].ListenCount
	})

	results := make([]client.SearchResult, 0)

	for _, a := range albums {
		results = append(results, a)
	}

	return results, nil
}

func NewMusicClient(httpClient *http.Client) *client.Client {
	return client.NewClient(
		httpClient,
		"https://musicbrainz.org/ws/2/release-group",
		"album",
		map[string]string{
			"fmt": "json",
		},
		map[string]string{
			"User-Agent": "test/1.0 ( dial1764@proton.me )",
		},
		readToSearchResult,
	)
}
