package geometry

import (
	"context"
	"log"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// DifferenceFeatureCollection returns a `geojson.FeatureCollection` representing the difference between the first (Feature) elements 'source_fc' and 'clip_fc'.
func DifferenceFeatureCollection(ctx context.Context, source_fc *geojson.FeatureCollection, clip_fc *geojson.FeatureCollection) (*geojson.FeatureCollection, error) {

	source_fc, err := UnionFeatureCollection(ctx, source_fc)

	if err != nil {
		return nil, err
	}

	clip_fc, err = UnionFeatureCollection(ctx, clip_fc)

	if err != nil {
		return nil, err
	}

	source_f := source_fc.Features[0]
	clip_f := clip_fc.Features[0]

	source_geom, err := OrbGeometryToGeosGeometry(ctx, source_f.Geometry)

	if err != nil {
		return nil, err
	}

	clip_geom, err := OrbGeometryToGeosGeometry(ctx, clip_f.Geometry)

	if err != nil {
		return nil, err
	}

	/*
		new_geom, err := source_geom.Difference(clip_geom)

		if err != nil {
			return nil, err
		}
	*/

	new_geom := source_geom.Difference(clip_geom)

	/*
		_, err = new_geom.Type()

		if err != nil {
			return nil, err
		}
	*/

	orb_geom, err := GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, err
	}

	if orb_geom.GeoJSONType() == "GeometryCollection" {

		log.Printf("WARNING difference-ed geometry returned as GeometryCollection\n")

		mp := make([]orb.Polygon, 0)

		for idx, g := range orb_geom.(orb.Collection) {

			t := g.GeoJSONType()
			switch t {
			case "Polygon":
				mp = append(mp, g.(orb.Polygon))
			case "MultiPolygon":

				for _, p := range g.(orb.MultiPolygon) {
					mp = append(mp, p)
				}
			default:
				log.Printf("WARNING geometry %d is a %s, skipping", idx, t)
			}
		}

		orb_geom = orb.MultiPolygon(mp)
	}

	new_f := geojson.NewFeature(orb_geom)
	new_f.Properties["foo"] = "bar"

	new_fc := geojson.NewFeatureCollection()
	new_fc.Append(new_f)

	return new_fc, nil
}
