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

func featureCollection(path string) (*geojson.FeatureCollection, error) {

	fh, err := os.Open(path)

	if err != nil {
		fmt.Errorf("Failed to open %s, %v", path, err)
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		fmt.Errorf("Failed to read %s, %v", path, err)
	}

	return geojson.UnmarshalFeatureCollection(body)
}

func main() {

	source_path := flag.String("source", "", "...")
	clip_path := flag.String("clip", "", "...")

	flag.Parse()

	ctx := context.Background()

	source_fc, err := featureCollection(*source_path)

	if err != nil {
		log.Fatalf("Failed to parse %s, %v", *source_path, err)
	}

	clip_fc, err := featureCollection(*clip_path)

	if err != nil {
		log.Fatalf("Failed to parse %s, %v", *clip_path, err)
	}

	new_fc, err := geometry.DifferenceFeatureCollection(ctx, source_fc, clip_fc)

	if err != nil {
		log.Fatalf("Failed to union %v", err)
	}

	enc, err := new_fc.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal union of %v", err)
	}

	fmt.Println(string(enc))
}
