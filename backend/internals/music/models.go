package music

type Response struct {
	Albums struct {
		Items []struct {
			Id      string `json:"id"`
			Title   string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"items"`
		Next string `json:"next"`
	} `json:"albums"`
}

type Album struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
	Next   string `json:"next"`
}
