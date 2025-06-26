package git

import (
	"context"
	"io"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func TestGitIterator(t *testing.T) {

	ctx := context.Background()

	it, err := iterate.NewIterator(ctx, "git:///tmp")

	if err != nil {
		t.Fatalf("Failed to create iterator, %v", err)
	}

	defer it.Close()

	iter_uri := "https://github.com/sfomuseum-data/sfomuseum-data-maps.git"

	for rec, err := range it.Iterate(ctx, iter_uri) {

		if err != nil {
			t.Fatalf("Failed to walk '%s', %v", iter_uri, err)
			break
		}

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
}
