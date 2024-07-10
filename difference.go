package geometry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/twpayne/go-geos"
)

// DifferenceFeatureCollection returns a `geojson.FeatureCollection` representing the difference between the first (Feature) elements 'source_fc' and 'clip_fc'.
func DifferenceFeatureCollection(ctx context.Context, source_fc *geojson.FeatureCollection, clip_fc *geojson.FeatureCollection) (*geojson.FeatureCollection, error) {

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

	new_geom := source_geom.Difference(clip_geom)

	orb_geom, err := GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive Orb geometry for new feature, %w", err)
	}

	if orb_geom.GeoJSONType() == "GeometryCollection" {

		slog.Warn("Difference-ed geometry returned as GeometryCollection\n")

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
				slog.Warn("Invalid geometry, skipping", "offset", idx, "type", t)
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

// DifferenceGeometriesWithFeatures returns a `geojson.Feature` instance representing the difference between 'base' and 'others'.
func DifferenceGeometriesWithFeatures(ctx context.Context, base *geojson.Feature, others ...*geojson.Feature) (*geojson.Feature, error) {

	base_geom, err := OrbGeometryToGeosGeometry(ctx, base.Geometry)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive GEOS geometry for base feature, %w", err)
	}

	other_geoms := make([]*geos.Geom, len(others))

	for idx, f := range others {

		geom, err := OrbGeometryToGeosGeometry(ctx, f.Geometry)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive GEOS geometry for other geometry at offset %d, %w", idx, err)
		}

		other_geoms[idx] = geom
	}

	new_geom, err := DifferenceGeometries(ctx, base_geom, other_geoms...)

	if err != nil {
		return nil, fmt.Errorf("Failed to diference geometries, %w", err)
	}

	orb_geom, err := GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive Orb geometry for new geometry, %w", err)
	}

	new_f := geojson.NewFeature(orb_geom)
	new_f.Properties = base.Properties

	return new_f, nil
}

// DifferenceGeometries returns a `geos.Geom` instance representing the difference between 'base' and 'others'.
func DifferenceGeometries(ctx context.Context, base_geom *geos.Geom, other_geoms ...*geos.Geom) (*geos.Geom, error) {

	to_remove := geos.NewEmptyPolygon()

	for _, g := range other_geoms {

		new_geom := to_remove.Union(g)
		to_remove = new_geom
	}

	new_geom := base_geom.Difference(to_remove)

	return new_geom, nil
}
