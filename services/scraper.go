package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapePage(url string) (string, error) {

	selectors := "main p, article p, section p, div p"

	httpClient := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %v %v", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	seen := make(map[string]struct{})
	var scrapedText []string

	doc.Find(selectors).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			if _, ok := seen[text]; !ok {
				seen[text] = struct{}{}
				scrapedText = append(scrapedText, text)
			}
		}
	})

	docText := strings.Join(scrapedText, "\n")

	return docText, nil
}
