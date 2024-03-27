package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type XMLpost struct {
	Channel struct {
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          struct {
			Title       string    `xml:"title"`
			Link        string    `xml:"link"`
			PubDate     time.Time `xml:"pubDate"`
			Guid        string    `xml:"guid"`
			Description string    `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type XmlPostResp struct {
	Item          struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		PubDate     time.Time `xml:"pubDate"`
		Guid        string    `xml:"guid"`
		Description string    `xml:"description"`
	} `xml:"item"`
}
// grab urls from the request and put them into a slice to handle in the
// enpoint handler LIKELY CASE XML
func (cfg *ApiConfig) getFeedsFromUrl(r *http.Response) ([]Post, error) {
	contentType := r.Header.Get("Content-Type")

	fmt.Printf("content type checked: %v", contentType)

	switch contentType {
	case "application/json":
		fmt.Println(": type json")

		type parameters struct {
			Url Post `json:"url"`
		}

		ScrapedUrls := make([]Post, 10)

		decode := json.NewDecoder(r.Body)
		params := parameters{}
		err := decode.Decode(&params)
		if err != nil {
			err = errors.New("cannot parse json")
			return nil, err
		}
		for i := 0; i <= 10; i++ {
			ScrapedUrls = append(ScrapedUrls, params.Url)
		}
		return ScrapedUrls, nil

	case "application/xml":
		fmt.Println(": type xml")

		type parameters struct {
			Url Post `xml:"link"`
		}

		ScrapedUrls := make([]Post, 10)

		decode := xml.NewDecoder(r.Body)
		params := parameters{}
		err := decode.Decode(&params)
		if err != nil {
			err = errors.New("cannot parse xml")
			return nil, err
		}
		for i := 0; i <= 10; i++ {
			ScrapedUrls = append(ScrapedUrls, params.Url)
		}
		return ScrapedUrls, nil
	default:
		fmt.Println("returning default")
		return nil, errors.New("cannot parse, bad request")
	}
}
