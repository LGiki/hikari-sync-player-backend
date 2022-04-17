package entity

type Episode struct {
	Title        string `json:"title"`
	PodcastName  string `json:"podcastName"`
	CoverUrl     string `json:"coverUrl"`
	EnclosureUrl string `json:"enclosureUrl"`
	ThemeColor   string `json:"themeColor"`
}
