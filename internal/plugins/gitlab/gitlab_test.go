package gitlabdoc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/apex/log"

	"github.com/stretchr/testify/require"
)

func TestRecentlyCreatedProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn(r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/user") {
			r, err := os.Open("testdata/user.json")
			require.NoError(t, err)
			_, err = io.Copy(w, r)
			require.NoError(t, err)
			return
		} else if strings.HasSuffix(r.URL.Path, "/projects") {
			r, err := os.Open("testdata/projects.json")
			require.NoError(t, err)
			_, err = io.Copy(w, r)
			require.NoError(t, err)
			return
		}
		defer r.Body.Close()
	}))
	defer srv.Close()
	os.Setenv("GITLAB_URL", srv.URL)

	p, err := recentlyCreatedProjects(5)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(p), 1)
}
