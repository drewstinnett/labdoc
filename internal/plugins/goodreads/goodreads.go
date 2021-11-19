package goodreads

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/mmcdole/gofeed"
)

type plug struct{}

// TemplateFuncMap returns a map of functions to be used in templates
func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"recentlyRead": recentlyRead,
	}
	return templ, nil
}

// Examples returns a list of example entries
func (p *plug) Examples() string {
	return `## Recently Read Books

{{ range goodreadsRecentlyRead 10 }}
- {{ .Title }} by {{ .Author }} {{ .PubDate | builtinAgo }}
{{- end }`
}

// recentlyRead returns a list of recently read books
func recentlyRead(limit int) ([]ReadEntry, error) {
	user := os.Getenv("GOODREADS_RSSUSERID")
	if user == "" {
		return nil, errors.New("Must Set GOODREADS_RSSUSERID")
	}
	key := os.Getenv("GOODREADS_RSSKEY")
	if key == "" {
		return nil, errors.New("Must Set GOODREADS_RSSKEY")
	}
	fp := gofeed.NewParser()
	url := fmt.Sprintf("https://www.goodreads.com/review/list_rss/%v?key=%v&shelf=read", user, key)
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	items := parseRSSRead(feed, limit)

	return items, nil
}

// ReadEntry is a single entry from the RSS feed
type ReadEntry struct {
	Title       string
	Link        string
	PubDate     *time.Time
	Description string
	Author      string
	UserRating  *float64
	// DateAdded   *time.Time
}

// parseRSSRead parses the RSS feed and returns a list of ReadEntry
func parseRSSRead(feed *gofeed.Feed, limit int) []ReadEntry {
	entries := make([]ReadEntry, 0, limit)

	for _, fi := range feed.Items {
		// Get the user rating if we can
		var userRating *float64
		userRatingGot, err := strconv.ParseFloat(fi.Custom["user_rating"], 64)
		if err != nil {
			userRating = nil
		} else {
			userRating = &userRatingGot
		}
		e := ReadEntry{
			Title:       fi.Title,
			Link:        fi.Link,
			PubDate:     fi.PublishedParsed,
			Description: fi.Description,
			Author:      fi.Custom["author_name"],
			UserRating:  userRating,
			// DateAdded:   fi.Custom["date_added"],
		}
		entries = append(entries, e)
		if len(entries) == limit {
			break
		}
	}

	return entries
}

func init() {
	labdoc.Add("goodreads", func() labdoc.Plugin {
		return &plug{}
	})
}
