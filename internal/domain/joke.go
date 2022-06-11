package domain

type JokeInfo struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Joke     string `json:"joke"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
	ID       int    `json:"id"`
}

type Jokes struct {
	List []*JokeInfo `json:"jokes"`
}
