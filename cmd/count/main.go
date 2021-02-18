package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

func main() {

	valid_schemes := strings.Join(emitter.Schemes(), ",")
	emitter_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported emitter URI schemes are: %s", valid_schemes)

	var emitter_uri = flag.String("emitter-uri", "git://", emitter_desc)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var count int64
	count = 0

	emitter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		_, err := emitter.PathForContext(ctx)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
		return nil
	}

	ctx := context.Background()

	iter, err := iterator.NewIterator(ctx, *emitter_uri, emitter_cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = iter.IterateURIs(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Counted %d records (saw %d records)\n", count, iter.Seen)
}
