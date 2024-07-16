package gpxtransform

import (
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/tkrajina/gpxgo/gpx"
)

/*
Filter creates a "GPXSegmentTransform"er that filters points from a segment based on filterOptions.
*/
func Filter(filterOptions ...options.FilterOptions) config.GPXSegmentTransform {
	return func(trackSegment gpx.GPXTrackSegment) (segments []gpx.GPXTrackSegment, err error) {
		points := []gpx.GPXPoint{}
		for index := 0; index < len(trackSegment.Points); index++ {
			doFilter := true
			for _, filterOption := range filterOptions {
				var doFilter_ bool
				doFilter_, err = filterOption.Do(trackSegment, index)
				if err != nil {
					return
				}
				if !doFilter_ {
					doFilter = false
					break
				}
			}
			if doFilter {
				points = append(points, trackSegment.Points[index])
			}
		}
		return []gpx.GPXTrackSegment{gpx.GPXTrackSegment{points, trackSegment.Extensions}}, err
	}
}
