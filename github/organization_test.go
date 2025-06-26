package github

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func TestIterateOrganization(t *testing.T) {

	ctx := context.Background()

	iter_uri := "githuborg://"

	expected := 43
	count := 0

	it, err := iterate.NewIterator(ctx, iter_uri)

	if err != nil {
		t.Fatalf("Failed to create iterator, %v", err)
	}

	defer it.Close()

	for rec, err := range it.Iterate(ctx, "sfomuseum-data://?prefix=sfomuseum-data-map") {

		if err != nil {
			t.Fatalf("Failed to iterate organization, %v", err)
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

		count += 1
	}

	if count != expected {
		t.Fatalf("Unexpected %d count but got: %d", expected, count)
	}
}
