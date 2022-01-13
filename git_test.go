package git

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"io"
	"testing"
)

func TestGitIterator(t *testing.T) {

	ctx := context.Background()

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		fmt.Println(path)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, "git:///tmp", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, "https://github.com/sfomuseum-data/sfomuseum-data-maps.git")

	if err != nil {
		t.Fatalf("Failed to iterate URIs, %v", err)
	}
}
