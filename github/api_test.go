package github

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"net/url"
	"testing"
	"time"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

var tests_access_token = flag.String("tests-access-token", "", "A valid GitHub API access token.")

func TestGitIterator(t *testing.T) {

	if *tests_access_token == "" {
		t.Skip()
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Verbose logging enabled")

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	it_q := url.Values{}
	it_q.Set("access_token", *tests_access_token)
	it_q.Set("concurrent", "true")

	it_url := url.URL{}
	it_url.Scheme = "githubapi"
	it_url.Host = "sfomuseum-data"
	it_url.Path = "/sfomuseum-data-maps"
	it_url.RawQuery = it_q.Encode()

	it, err := iterate.NewIterator(ctx, it_url.String())

	if err != nil {
		t.Fatalf("Failed to create iterator, %v", err)
	}

	it_path := "data"

	for rec, err := range it.Iterate(ctx, it_path) {

		if err != nil {
			t.Fatalf("Failed to walk %s, %v", it_path, err)
			break
		}

		slog.Debug("Process record", "path", rec.Path)

		defer rec.Body.Close()
		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed to read body for %s, %v", rec.Path, err)
		}

		_, err = rec.Body.Seek(0, 0)

		if err != nil {
			t.Fatalf("Failed to rewind body for %s, %v", rec.Path, err)
		}

		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed second read body for %s, %v", rec.Path, err)
		}
	}

	seen := it.Seen()
	expected := int64(43)

	if seen != expected {
		t.Fatalf("Unexpected record count. Got %d but expected %d", seen, expected)
	}

	err = it.Close()

	if err != nil {
		t.Fatalf("Failed to close iterator, %v", err)
	}
}
