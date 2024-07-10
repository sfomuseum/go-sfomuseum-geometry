// package geometry provides method for geometric operations related to SFO Museum (Who's On First) records.
package geometry

import (
	"bytes"
	"context"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/twpayne/go-geos"
)

// OrbGeometryToGeosGeometry converts an `orb.Geometry` instance in to a `geos.Geom` instance.
func OrbGeometryToGeosGeometry(ctx context.Context, orb_geom orb.Geometry) (*geos.Geom, error) {
	str_wkt := wkt.MarshalString(orb_geom)
	return geos.NewGeomFromWKT(str_wkt)
}

// GeosGeometryToOrbGeometry converts a `geos.Geom`	instance in to an `orb.Geometry` instance.
func GeosGeometryToOrbGeometry(ctx context.Context, geos_geom *geos.Geom) (orb.Geometry, error) {
	wkb_body := geos_geom.ToWKB()
	br := bytes.NewReader(wkb_body)
	dec := wkb.NewDecoder(br)
	return dec.Decode()
}
