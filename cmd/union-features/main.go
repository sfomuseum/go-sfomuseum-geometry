package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-sfomuseum-geometry"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	paths := flag.Args()

	features := make([]*geojson.Feature, 0)

	for _, path := range paths {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s, %v", path, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			log.Fatalf("Failed to read %s, %v", path, err)
		}

		fc, err := geojson.UnmarshalFeature(body)

		if err != nil {
			log.Fatalf("Failed to parse %s, %v", path, err)
		}

		features = append(features, fc)
	}

	new_f, err := geometry.UnionFeatures(ctx, features...)

	if err != nil {
		log.Fatalf("Failed to union FeatureCollections, %v", err)
	}

	enc, err := new_f.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal union of features, %v", err)
	}

	fmt.Println(string(enc))
}
