package movies

type Response struct {
	Results []struct {
		Id         int     `json:"id"`
		Title      string  `json:"title"`
		Popularity float32 `json:"popularity"`
		Poster     string  `json:"poster_path"`
	} `json:"results"`
}

type Movie struct {
	Id         int     `json:"id"`
	Title      string  `json:"title"`
	Popularity float32 `json:"popularity"`
	Poster     string  `json:"poster"`
}
