package options

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

/*
DirectOptions hold a function that directly modifies a segment,
returning zero, one or multiple modified segments.
Unlike Split and Filter, Direct functions are able to apply any type of modification.
*/
type DirectOptions struct {
	Do func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error)
}

/*
Polyline simplification based on Ramer-Douglas-Peucker:
Points are removed when lying within maxDistance of the line formed by its two
adjacent points.
*/
func DouglasPeucker(maxDistance unit.Length) DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			segment.SimplifyTracks(float64(maxDistance))
			return []gpx.GPXTrackSegment{segment}, nil
		},
	}
}

/*
RemoveStops removes points that are assumed to have no relevant movement.
See PauseSplit for a function that splits segments, if a pause is detected.
*/
func RemoveStops() DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			filteredPoints := []gpx.GPXPoint{}
			trackPositions := segment.StoppedPositions()
			trackIndex := 0
			for i, point := range segment.Points {
				if trackIndex < len(trackPositions) && trackPositions[trackIndex].PointNo == i {
					trackIndex++
					slog.Info(fmt.Sprintf("RemoveStops: Point %v captured at %v is removed", i, point.Timestamp))
				} else {
					filteredPoints = append(filteredPoints, point)
				}
			}
			return []gpx.GPXTrackSegment{gpx.GPXTrackSegment{Points: filteredPoints, Extensions: segment.Extensions}}, nil
		},
	}
}

/*
MinPoints removes segments with less points than the provided minPoints
*/
func MinPoints(minPoints int) DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			if len(segment.Points) < minPoints {
				return []gpx.GPXTrackSegment{}, nil
			}
			return []gpx.GPXTrackSegment{segment}, nil
		},
	}
}

/*
MinDistance removes segments with a shorter distance than the provided minDistance
*/
func MinDistance(minDistance unit.Length) DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			if len(segment.Points) == 0 {
				return []gpx.GPXTrackSegment{segment}, nil
			}
			if float64(minDistance) <= segment.Length3D() {
				return []gpx.GPXTrackSegment{segment}, nil
			}
			return []gpx.GPXTrackSegment{}, nil
		},
	}
}

/*
MinRadius removes segments, if all their points lie within a given minRadius from the starting point
*/
func MinRadius(minRadius unit.Length) DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			for _, point := range segment.Points {
				if float64(minRadius) < gpx.Length3D([]gpx.Point{segment.Points[0].Point, point.Point}) {
					return []gpx.GPXTrackSegment{segment}, nil
				}
			}
			return []gpx.GPXTrackSegment{}, nil
		},
	}
}

/*
MinDuration removes segments with a shorter duration than the provided minDuration
*/
func MinDuration(minDuration time.Duration) DirectOptions {
	remove := func(actualDuration time.Duration) bool {
		return actualDuration < minDuration
	}
	return segmentDurationFilter(remove)
}

/*
MaxDuration removes segments with a longer duration than the provided maxDuration
*/
func MaxDuration(maxDuration time.Duration) DirectOptions {
	remove := func(actualDuration time.Duration) bool {
		return actualDuration > maxDuration
	}
	return segmentDurationFilter(remove)
}

/*
segmentDurationFilter removes segments, if the remove function returns true when passing their corresponding segment duration.
*/
func segmentDurationFilter(remove func(actualDuration time.Duration) bool) DirectOptions {
	return DirectOptions{
		func(segment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error) {
			if len(segment.Points) == 0 {
				return []gpx.GPXTrackSegment{segment}, nil
			}
			duration := segment.Points[len(segment.Points)-1].Timestamp.Sub(segment.Points[0].Timestamp)
			if remove(duration) {
				return []gpx.GPXTrackSegment{}, nil
			}
			return []gpx.GPXTrackSegment{segment}, nil
		},
	}
}
