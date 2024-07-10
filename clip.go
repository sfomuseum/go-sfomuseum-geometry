package geometry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// ClipFeatureCollection returns a `geojson.FeatureCollection` instance representing the intersection of 'source_fc' and 'clip_fc'.
func ClipFeatureCollection(ctx context.Context, source_fc *geojson.FeatureCollection, clip_fc *geojson.FeatureCollection) (*geojson.FeatureCollection, error) {

	source_fc, err := UnionFeatureCollection(ctx, source_fc)

	if err != nil {
		return nil, fmt.Errorf("Failed to union source feature collection, %w", err)
	}

	clip_fc, err = UnionFeatureCollection(ctx, clip_fc)

	if err != nil {
		return nil, fmt.Errorf("Failed to union clip feature collection, %w", err)
	}

	source_f := source_fc.Features[0]
	clip_f := clip_fc.Features[0]

	source_geom, err := OrbGeometryToGeosGeometry(ctx, source_f.Geometry)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive GEOS geometry for source feature (offset 0), %w", err)
	}

	clip_geom, err := OrbGeometryToGeosGeometry(ctx, clip_f.Geometry)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive GEOS geometry for clip feature (offset 0), %w", err)
	}

	new_geom := source_geom.Intersection(clip_geom)

	// t := new_geom.Type()

	orb_geom, err := GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive Orb geometry for new geometry, %w", err)
	}

	if orb_geom.GeoJSONType() == "GeometryCollection" {

		slog.Warn("Clipped geometry returned as GeometryCollection\n")

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
				slog.Warn("Invalid geometry", "offset", idx, "type", t)
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
