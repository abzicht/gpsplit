package options

import (
	"time"

	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

/*
SplitOptions hold a function that returns true iff a segment should be split at an indicated point
*/
type SplitOptions TransformOptions

/*
A time-based splitter that splits whenever two consecutive points have a time
difference larger than maxStep.
*/
func TimeSplit(maxStep time.Duration) SplitOptions {
	return SplitOptions{
		func(segment gpx.GPXTrackSegment, index int) (bool, error) {
			duration := segment.Points[index+1].Timestamp.Sub(segment.Points[index].Timestamp)
			return maxStep <= duration, nil
		},
	}
}

/*
A distance-based splitter that splits whenever two consecutive points have a
distance larger than maxDist.
*/
func DistanceSplit(maxDist unit.Length) SplitOptions {
	return SplitOptions{
		func(segment gpx.GPXTrackSegment, index int) (bool, error) {
			dist := gpx.Length3D([]gpx.Point{segment.Points[index].Point, segment.Points[index+1].Point})
			return dist > float64(maxDist), nil
		},
	}
}

/*
PauseSplit splits when previous points lie within the provided radius (from the currently viewed point) for at least minDuration. It waits with splitting until the next point (from the currently viewed point) would lie outside the radius.
Different to RemoveStops, this function can be configured to not affect shorter pauses.
*/
func PauseSplit(radius unit.Length, minDuration time.Duration) SplitOptions {
	return SplitOptions{
		func(segment gpx.GPXTrackSegment, index int) (bool, error) {
			point := segment.Points[index]
			for i := index - 1; i >= 0; i-- {
				prevPoint := segment.Points[i]
				tooLongAgo := point.Timestamp.Sub(prevPoint.Timestamp) > minDuration
				tooFarAway := gpx.Length3D([]gpx.Point{prevPoint.Point, point.Point}) > float64(radius)
				if tooLongAgo {
					if tooFarAway {
						// edge case: we moved sufficiently out of the radius
						// just before we would have split
						return false, nil
					} else {
						// the last index-i points are within the radius and we
						// exceeded the minDuration. We split between the index
						// point and the next one. But we only do that, if the
						// next point is out of the radius! So let's check
						// that.
						if index >= len(segment.Points)-1 {
							return false, nil
						}
						if gpx.Length3D([]gpx.Point{point.Point, segment.Points[index+1].Point}) > float64(radius) {
							// the next point is too far away, it is time to
							// split!
							return true, nil
						} else {
							return false, nil
						}
					}
				} else {
					if tooFarAway {
						// normal case: we moved sufficiently out of the radius within the
						// minDuration
						return false, nil
					} else {
						// this point lies within the radius and is quite
						// recent - we may still be traveling and therefore do
						// nothing.
						continue
					}
				}
			}
			return false, nil
		},
	}
}
