package scraper

import (
	"bytes"

	"github.com/gocolly/colly"
)

func Scrape(url string) (string, error) {
	var buffer bytes.Buffer
	c := colly.NewCollector()

	c.OnHTML("p, h1, h2, h3, h4, h5, h6, div, span, li, td", func(e *colly.HTMLElement) {
		buffer.WriteString(e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	})

	return buffer.String(), nil
}
