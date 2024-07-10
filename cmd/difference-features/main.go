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

func loadFeature(path string) (*geojson.Feature, error) {

	r, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", path, err)
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read %s, %w", path, err)
	}

	return geojson.UnmarshalFeature(body)
}

func main() {

	var source_path string

	flag.StringVar(&source_path, "source", "", "...")

	flag.Parse()

	other_paths := flag.Args()

	ctx := context.Background()

	source_f, err := loadFeature(source_path)

	if err != nil {
		log.Fatalf("Failed to load feature %s, %v", source_path, err)
	}

	other_f := make([]*geojson.Feature, len(other_paths))

	for idx, path := range other_paths {

		f, err := loadFeature(path)

		if err != nil {
			log.Fatalf("Failed to load feature %s, %v", path, err)
		}

		other_f[idx] = f
	}

	new_f, err := geometry.DifferenceGeometriesWithFeatures(ctx, source_f, other_f...)

	if err != nil {
		log.Fatalf("Failed to difference geometries %v", err)
	}

	enc, err := new_f.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal union of %v", err)
	}

	fmt.Println(string(enc))
}
