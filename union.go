package geometry

import (
	"context"
	"github.com/paulmach/orb/geojson"
	"github.com/paulsmith/gogeos/geos"
	"log"
)

// debug

// UnionFeatureCollection returns a `geojson.FeatureCollection` instance with a single feature representing the union
// of all the features in 'cols'
func UnionFeatureCollection(ctx context.Context, cols ...*geojson.FeatureCollection) (*geojson.FeatureCollection, error) {

	var features []*geojson.Feature

	if len(cols) == 1 {
		features = cols[0].Features
	} else {

		features = make([]*geojson.Feature, 0)

		for _, fc := range cols {

			for _, f := range fc.Features {

				features = append(features, f)
			}
		}
	}

	new_fc := geojson.NewFeatureCollection()

	first := features[0]
	orb_geom := first.Geometry

	new_geom, err := OrbGeometryToGeosGeometry(ctx, orb_geom)

	if err != nil {
		return nil, err
	}

	for idx, f := range features[1:] {

		geos_g, err := OrbGeometryToGeosGeometry(ctx, f.Geometry)

		if err != nil {
			return nil, err
		}

		g, err := new_geom.Union(geos_g)

		if err != nil {
			return nil, err
		}

		t, err := g.Type()

		if err != nil {
			return nil, err
		}

		// But why? (20210329/thisisaaronland)

		if t == geos.GEOMETRYCOLLECTION {
			log.Printf("WARNING feature %d causes union to produce a GeometryCollection, skipping\n", idx)
			continue
		}

		new_geom = g
	}

	orb_geom, err = GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, err
	}

	new_props := map[string]interface{}{
		"hello": "world",
	}

	new_f := geojson.NewFeature(orb_geom)
	new_f.Properties = new_props

	new_fc.Append(new_f)

	return new_fc, nil
}
