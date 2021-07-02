// package geometry provides method for geometric operations related to SFO Museum (Who's On First) records.
package geometry

import (
	"bytes"
	"context"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulsmith/gogeos/geos"
)

func OrbGeometryToGeosGeometry(ctx context.Context, orb_geom orb.Geometry) (*geos.Geometry, error) {
	str_wkt := wkt.MarshalString(orb_geom)
	return geos.FromWKT(str_wkt)
}

func GeosGeometryToOrbGeometry(ctx context.Context, geos_geom *geos.Geometry) (orb.Geometry, error) {

	wkb_body, err := geos_geom.WKB()

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(wkb_body)

	dec := wkb.NewDecoder(br)
	return dec.Decode()
}

func DifferenceGeometriesWithFeatures(ctx context.Context, base *geojson.Feature, others ...*geojson.Feature) (*geojson.Feature, error) {

	base_geom, err := OrbGeometryToGeosGeometry(ctx, base.Geometry)

	if err != nil {
		return nil, err
	}

	other_geoms := make([]*geos.Geometry, len(others))

	for idx, f := range others {

		geom, err := OrbGeometryToGeosGeometry(ctx, f.Geometry)

		if err != nil {
			return nil, err
		}

		other_geoms[idx] = geom
	}

	new_geom, err := DifferenceGeometries(ctx, base_geom, other_geoms...)

	if err != nil {
		return nil, err
	}

	orb_geom, err := GeosGeometryToOrbGeometry(ctx, new_geom)

	if err != nil {
		return nil, err
	}

	new_f := geojson.NewFeature(orb_geom)
	new_f.Properties = base.Properties

	return new_f, nil
}

func DifferenceGeometries(ctx context.Context, base_geom *geos.Geometry, other_geoms ...*geos.Geometry) (*geos.Geometry, error) {

	to_remove, err := geos.EmptyPolygon()

	if err != nil {
		return nil, err
	}

	for _, g := range other_geoms {

		new_geom, err := to_remove.Union(g)

		if err != nil {
			return nil, err
		}

		to_remove = new_geom
	}

	new_geom, err := base_geom.Difference(to_remove)

	if err != nil {
		return nil, err
	}

	return new_geom, nil
}
