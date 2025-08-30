package options

import (
	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

/*
FilterOptions hold a function that returns true iff a point should be included in the filter result
*/
type FilterOptions TransformOptions

/*
TrimStart cuts all points that are reached from the first recorded point within the provided max radius
*/
func TrimStart(maxRadius unit.Length) FilterOptions {
	return FilterOptions{
		func(segment gpx.GPXTrackSegment, index int) (bool, error) {
			if len(segment.Points) == 0 {
				return false, nil
			}
			if index == 0 {
				return true, nil
			}

			for i, point := range segment.Points {
				if float64(maxRadius) < point.Distance3D(&segment.Points[0]) {
					break
				}
				if i == index {
					return false, nil
				}
			}
			return true, nil
		},
	}
}

/*
TrimEnd cuts all points that are reached backwards from the last recorded point within the provided max radius
*/
func TrimEnd(maxRadius unit.Length) FilterOptions {
	return FilterOptions{
		func(segment gpx.GPXTrackSegment, index int) (bool, error) {
			if len(segment.Points) == 0 {
				return false, nil
			}
			if index == len(segment.Points)-1 {
				return true, nil
			}

			for i := len(segment.Points) - 1; i >= 0; i-- {
				if float64(maxRadius) < segment.Points[i].Distance3D(&segment.Points[len(segment.Points)-1]) {
					break
				}
				if i == index {
					return false, nil
				}
			}
			return true, nil
		},
	}
}
