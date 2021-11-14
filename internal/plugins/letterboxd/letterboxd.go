package letterboxd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/mmcdole/gofeed"
)

type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"recentlyWatched": recentlyWatched,
	}
	return templ, nil
}

func (p *plug) Examples() string {
	return `## Watched List
{{ range letterboxdRecentlyWatched 5 }}
* {{ .Verb }} [{{ .TitleWithRating }})]({{ .Link }})
{{- end }}
	
## Watched Posters
{{ range letterboxdRecentlyWatched 5 }}
![{{ .FilmTitle }}]({{ .Poster}})
{{- end }}
`
}

type WatchedEntry struct {
	FilmTitle       string
	TitleWithRating string
	FilmYear        int
	Description     string
	Creator         string
	WatchedDate     time.Time
	PubDate         time.Time
	MemberRating    float64
	Link            string
	Rewatch         bool
	Poster          string
	Verb            string
}

func recentlyWatched(limit int) ([]WatchedEntry, error) {
	user := os.Getenv("LETTERBOXD_USER")
	if user == "" {
		return nil, errors.New("Must Set LETTERBOXD_USER")
	}
	fp := gofeed.NewParser()
	url := fmt.Sprintf("https://letterboxd.com/%v/rss/", user)
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	items := parseRSSWatched(feed, limit)
	return items, nil
}

func parseRSSWatched(feed *gofeed.Feed, limit int) []WatchedEntry {
	var entries []WatchedEntry
	for _, fi := range feed.Items {
		// Only return movies with a watchedDate set
		if _, ok := fi.Extensions["letterboxd"]["watchedDate"]; ok {
			year := fi.Extensions["letterboxd"]["filmYear"][0].Value
			yearI, _ := strconv.Atoi(year)

			memberRating := fi.Extensions["letterboxd"]["memberRating"][0].Value
			memberRatingF, _ := strconv.ParseFloat(memberRating, 64)

			var rewatch bool
			if fi.Extensions["letterboxd"]["rewatch"][0].Value == "Yes" {
				rewatch = true
			}

			watchedDate := fi.Extensions["letterboxd"]["watchedDate"][0].Value
			watchedDateTime, _ := time.Parse("2006-01-02", watchedDate)

			e := WatchedEntry{
				FilmTitle:       fi.Extensions["letterboxd"]["filmTitle"][0].Value,
				TitleWithRating: fi.Title,
				FilmYear:        yearI,
				Creator:         fi.Extensions["dc"]["creator"][0].Value,
				MemberRating:    memberRatingF,
				Link:            fi.Link,
				Rewatch:         rewatch,
				WatchedDate:     watchedDateTime,
				PubDate:         *fi.PublishedParsed,
				Description:     fi.Description,
			}

			if rewatch {
				e.Verb = "Rewatched"
			} else {
				e.Verb = "Watched"
			}

			// Get the poster URL if we can
			poster, err := extractPoster(fi.Description)
			if err == nil {
				e.Poster = poster
			}

			entries = append(entries, e)
			if len(entries) == limit {
				break
			}
		}
	}
	return entries
}

func extractPoster(data string) (string, error) {
	re := regexp.MustCompile(`.*(https://.*resized.*)"\/>.*`)
	// re := regexp.MustCompile(`(.*)`)
	rs := re.FindStringSubmatch(data)
	if len(rs) < 2 {
		return "", errors.New("NoMatches")
	}
	return rs[1], nil
}

func init() {
	labdoc.Add("letterboxd", func() labdoc.Plugin {
		return &plug{}
	})
}
