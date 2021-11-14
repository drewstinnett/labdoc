package goodreads

import (
	"os"
	"testing"

	"github.com/apex/log"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func TestParseRecentlyRead(t *testing.T) {
	file, err := os.Open("./testdata/read.xml")
	require.NoError(t, err)
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	entries := parseRSSRead(feed, 10)
	require.Greater(t, len(entries), 0)
	first := entries[0]
	log.Warnf("%+v", first)

	require.Equal(t, "Blankets", first.Title)
	require.Equal(t, "Craig Thompson", first.Author)
	require.Equal(t, float64(5), *first.UserRating)
	require.Equal(t, "https://www.goodreads.com/review/show/3262380767?utm_medium=api&utm_source=rss", first.Link)
}
