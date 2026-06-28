package music

type Response struct {
	Results []struct {
		Id      string `json:"id"`
		Title   string `json:"title"`
		Credits []struct {
			Name string `json:"name"`
		} `json:"artist-credit"`
	} `json:"release-groups"`
}

type Album struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
}
