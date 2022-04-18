package xiaoyuzhou_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"hikari_sync_player/entity"
	"hikari_sync_player/pkg/request"
	"io/ioutil"
	"regexp"
)

type AssociatedMedia struct {
	Type       string `json:"@type"`
	ContentURL string `json:"contentUrl"`
}

type PartOfSeries struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type podcastShow struct {
	Context         string          `json:"@context"`
	Type            string          `json:"@type"`
	AssociatedMedia AssociatedMedia `json:"associatedMedia"`
	DatePublished   string          `json:"datePublished"`
	Description     string          `json:"description"`
	Name            string          `json:"name"`
	PartOfSeries    PartOfSeries    `json:"partOfSeries"`
	TimeRequired    string          `json:"timeRequired"`
	URL             string          `json:"url"`
}

func parseEpisodeFromHtml(html string) (*entity.PodcastEpisode, error) {
	podcastShowPattern := regexp.MustCompile(`(?m)<script name="schema:podcast-show" type="application/ld\+json">(?P<json>.*?)</script>`)
	coverImagePattern := regexp.MustCompile(`(?m)<meta property="og:image" content="(?P<url>.*?)"\/>`)
	themeColorPattern := regexp.MustCompile(`(?m)--theme-color: hsl\((?P<hsl>.*?)\);`)
	podcastShowMatchResult := podcastShowPattern.FindStringSubmatch(html)

	if len(podcastShowMatchResult) == 0 {
		return nil, errors.New("fail to parse xiaoyuzhou episode html")
	}
	var podcastShow podcastShow
	err := json.Unmarshal([]byte(podcastShowMatchResult[podcastShowPattern.SubexpIndex("json")]), &podcastShow)
	if err != nil {
		return nil, errors.New("fail to parse xiaoyuzhou episode json")
	}
	episode := &entity.PodcastEpisode{
		Title:        podcastShow.Name,
		PodcastName:  podcastShow.PartOfSeries.Name,
		EnclosureUrl: podcastShow.AssociatedMedia.ContentURL,
	}

	coverImageMatchResult := coverImagePattern.FindStringSubmatch(html)
	if len(coverImageMatchResult) != 0 {
		episode.CoverUrl = coverImageMatchResult[coverImagePattern.SubexpIndex("url")]
	}

	themeColorMatchResult := themeColorPattern.FindStringSubmatch(html)
	if len(themeColorMatchResult) != 0 {
		episode.ThemeColor = fmt.Sprintf("hsl(%s)", themeColorMatchResult[themeColorPattern.SubexpIndex("hsl")])
	}

	return episode, nil
}

func ParseEpisode(url string) (*entity.PodcastEpisode, error) {
	response, err := request.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	responseText := string(responseBytes)
	return parseEpisodeFromHtml(responseText)
}
