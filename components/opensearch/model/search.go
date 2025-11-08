package model

type SearchResultResp struct {
	ScrollId string               `json:"_scroll_id"`
	Took     int                  `json:"took"`
	TimedOut bool                 `json:"timed_out"`
	Shards   SearchResponseShards `json:"_shards"`
	Hits     SearchResponseHits   `json:"hits"`
}

type SearchDocItem struct {
	Index  string                 `json:"_index"`
	Id     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type SearchResponseShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type SearchResponseHits struct {
	Total    SearchResponseHitsTotal `json:"total"`
	MaxScore float64                 `json:"max_score"`
	Hits     []SearchDocItem         `json:"hits"`
}

type SearchResponseHitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
