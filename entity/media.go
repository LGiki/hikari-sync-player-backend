package entity

type PodcastEpisode struct {
	Title        string `json:"title"`
	PodcastName  string `json:"podcastName"`
	CoverUrl     string `json:"coverUrl"`
	EnclosureUrl string `json:"enclosureUrl"`
	ThemeColor   string `json:"themeColor"`
}

type Video struct {
	Url string `json:"title"`
}
