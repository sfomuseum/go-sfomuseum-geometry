package geometry

import (
	"context"

	"github.com/paulmach/orb/geojson"
)

// ContainedBy returns the list of `geojson.Feature` instances from 'candidates' that contain 'f'.
func ContainedBy(ctx context.Context, f *geojson.Feature, candidates ...*geojson.Feature) ([]*geojson.Feature, error) {

	contained_by := make([]*geojson.Feature, 0)

	for _, candidate_f := range candidates {

		is_contained, err := isContainedBy(ctx, f, candidate_f)

		if err != nil {
			return nil, err
		}

		if is_contained {
			contained_by = append(contained_by, candidate_f)
		}
	}

	return contained_by, nil
}

// IsContainedBy returns a boolean flag indicating whether 'f2' contains 'f1'.
func isContainedBy(ctx context.Context, f1 *geojson.Feature, f2 *geojson.Feature) (bool, error) {

	f1_geom, err := OrbGeometryToGeosGeometry(ctx, f1.Geometry)

	if err != nil {
		return false, err
	}

	f2_geom, err := OrbGeometryToGeosGeometry(ctx, f2.Geometry)

	if err != nil {
		return false, err
	}

	// Contains returns true if every point of the other is a point of this geometry, and the interiors of the two geometries have at least one point in common.

	return f2_geom.Contains(f1_geom), nil
}
