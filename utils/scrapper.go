package utils

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"github.com/etrinque/rssagg/database"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type XMLpost struct {
	Channel struct {
		Title       string        `xml:"title"`
		Link        string        `xml:"link"`
		Description string        `xml:"description"`
		Language    string        `xml:"language"`
		Item        []XmlPostResp `xml:"item"`
	} `xml:"channel"`
}

type XmlPostResp struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func startScraper(db *database.Queries, concurrency int, timeBetweenScrapes time.Duration) {
	tick := time.NewTicker(timeBetweenScrapes)
	for ; ; <-tick.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("could not get feeds", err)
			continue
		}

		log.Printf("found %d feeds.", len(feeds))
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeeds(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeeds(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("could not fetch feed from feed: %s", feed.Name)
		return
	}
	feedData, err := getFeedsFromUrl(feed.Url)
	if err != nil {
		log.Printf("could not get feeds from: %s", feed.Url)
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique contraint") {
				continue
			}
			log.Printf("Couldnt create post: %v", err)
			continue
		}
		log.Println("found feed", item.Title)
	}
	log.Printf("Feed %s retrieved, %v posts found", feed.Name, len(feedData.Channel.Item))

}

// grab urls from the request and put them into a slice to handle in the
// enpoint handler LIKELY CASE XML
func getFeedsFromUrl(feedURL string) (*XMLpost, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("logging resp: %v", resp)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not retrieve body")

	}

	log.Printf("logging data: %v", data)

	var xmlresp XMLpost
	err = xml.Unmarshal(data, &xmlresp)
	if err != nil {
		return nil, errors.New("error unmarshalling")
	}

	log.Printf("logging unmarshal: %v", &xmlresp)

	return &xmlresp, nil
}
