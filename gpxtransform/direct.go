package gpxtransform

import (
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/tkrajina/gpxgo/gpx"
)

func Direct(directOptions ...options.DirectOptions) config.GPXSegmentTransform {
	return func(trackSegment gpx.GPXTrackSegment) (segments []gpx.GPXTrackSegment, err error) {
		segments = []gpx.GPXTrackSegment{trackSegment}
		for _, directOption := range directOptions {
			for segmentIndex, _ := range segments {
				var newSegments []gpx.GPXTrackSegment
				newSegments, err = directOption.Do(segments[segmentIndex])
				if err != nil {
					return
				}
				if len(newSegments) == 0 {
					segments = append(segments[0:segmentIndex], segments[segmentIndex+1:len(segments)]...)
					continue
				}
				segments = append(segments[0:segmentIndex], append(newSegments, segments[segmentIndex+1:len(segments)]...)...)
			}
		}
		return
	}
}
