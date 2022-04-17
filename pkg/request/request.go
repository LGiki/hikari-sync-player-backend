package request

import (
	"net/http"
	"time"
)

func Get(url string) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
