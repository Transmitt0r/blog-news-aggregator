package newsapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	return http.DefaultClient.Do(req)
}

func (c *Client) GetTopHeadlines(q string) ([]Article, error) {
	url := "https://newsapi.org/v2/top-headlines"
	params := fmt.Sprintf("?q=%s", q)
	resp, err := c.doRequest(fmt.Sprintf("%s%s", url, params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Articles []Article `json:"articles"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Articles, nil
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
