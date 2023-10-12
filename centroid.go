package geometry

import (
	"context"

	"github.com/paulmach/orb/geojson"
)

// CentroidWithFeature return the centroid for 'f' as determined by GEOS.
func CentroidWithFeature(ctx context.Context, f *geojson.Feature) (float64, float64, error) {

	geos_geom, err := OrbGeometryToGeosGeometry(ctx, f.Geometry)

	if err != nil {
		return 0.0, 0.0, err
	}

	// What we really want is this...
	// https://github.com/mbloch/mapshaper/blob/4a1eac2845420472bb23df863723aa8e3021ced2/src/points/mapshaper-anchor-points.js

	/*
		centroid, err := geos_geom.Centroid()

		if err != nil {
			return 0.0, 0.0, err
		}

		lon, err := centroid.X()

		if err != nil {
			return 0.0, 0.0, err
		}

		lat, err := centroid.Y()

		if err != nil {
			return 0.0, 0.0, err
		}

	*/

	centroid := geos_geom.Centroid()

	lon := centroid.X()
	lat := centroid.Y()

	return lon, lat, nil
}
