package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-sfomuseum-geometry"
	"io"
	"log"
	"os"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	paths := flag.Args()

	cols := make([]*geojson.FeatureCollection, 0)

	for _, path := range paths {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s, %v", path, err)
		}

		defer fh.Close()

		body, err := io.ReadAll(fh)

		if err != nil {
			log.Fatalf("Failed to read %s, %v", path, err)
		}

		fc, err := geojson.UnmarshalFeatureCollection(body)

		if err != nil {
			log.Fatalf("Failed to parse %s, %v", path, err)
		}

		cols = append(cols, fc)
	}

	new_fc, err := geometry.UnionFeatureCollection(ctx, cols...)

	if err != nil {
		log.Fatalf("Failed to union FeatureCollections, %v", err)
	}

	enc, err := new_fc.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal union of FeatureCollections, %v", err)
	}

	fmt.Println(string(enc))
}
