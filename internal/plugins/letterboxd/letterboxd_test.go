package letterboxd

import (
	"os"
	"testing"

	"github.com/apex/log"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func TestParseRecentlyWatched(t *testing.T) {
	file, err := os.Open("./testdata/rss.xml")
	defer file.Close()
	require.NoError(t, err)
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	entries := parseRSSWatched(feed, 10)
	require.Greater(t, len(entries), 0)
	first := entries[0]
	log.Warnf("%+v", first)

	require.Equal(t, "Great White", first.FilmTitle)
	require.Equal(t, 2021, first.FilmYear)
	require.Equal(t, "Drew Stinnett", first.Creator)
	require.Equal(t, 2.5, first.MemberRating)
	require.Equal(t, "https://letterboxd.com/mondodrew/film/great-white-2021/", first.Link)
	require.Equal(t, false, first.Rewatch)
	require.Equal(t, ` <p><img src="https://a.ltrbxd.com/resized/film-poster/4/6/3/3/6/6/463366-great-white-0-500-0-750-crop.jpg?k=5af2dbce5b"/></p> <p>Watched on Thursday November 11, 2021.</p> `, first.Description)
	require.Equal(t, "https://a.ltrbxd.com/resized/film-poster/4/6/3/3/6/6/463366-great-white-0-500-0-750-crop.jpg?k=5af2dbce5b", first.Poster)
}

func TestExtractPoster(t *testing.T) {
	tests := []struct {
		data    string
		want    string
		wantErr bool
	}{
		{
			` <p><img src="https://a.ltrbxd.com/resized/film-poster/4/6/3/3/6/6/463366-great-white-0-500-0-750-crop.jpg?k=5af2dbce5b"/></p> <p>Watched on Thursday November 11, 2021.</p> `,
			"https://a.ltrbxd.com/resized/film-poster/4/6/3/3/6/6/463366-great-white-0-500-0-750-crop.jpg?k=5af2dbce5b",
			false,
		},
		{
			` <p><img src="https://a.ltrbxd.com/resized/sm/upload/pj/p0/i3/y9/aNxU4eZrfiChOgCEtyKWSCExdEi-0-500-0-750-crop.jpg?k=866b9f1460"/></p> <p>Watched on Thursday November 11, 2021.</p> `,
			"https://a.ltrbxd.com/resized/sm/upload/pj/p0/i3/y9/aNxU4eZrfiChOgCEtyKWSCExdEi-0-500-0-750-crop.jpg?k=866b9f1460",
			false,
		},
		{
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		got, err := extractPoster(tt.data)
		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tt.want, got)
	}
}
