package gpxtransform

import (
	"fmt"

	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/tkrajina/gpxgo/gpx"
)

func AnalyzeFile() config.GPXFileTransform {
	return func(gpxFile gpx.GPX) (files []gpx.GPX, err error) {
		fmt.Println(gpxFile.GetGpxInfo())
		return []gpx.GPX{}, nil
	}
}

/*
CountFile prints counts for different gpx objects
*/
func CountFile() config.GPXFileTransform {
	return func(gpxFile gpx.GPX) (files []gpx.GPX, err error) {
		numSegments := 0
		numPoints := 0
		for _, track := range gpxFile.Tracks {
			numSegments += len(track.Segments)
			for _, segment := range track.Segments {
				numPoints += len(segment.Points)
			}
		}
		avgNumSegments := numSegments / len(gpxFile.Tracks)
		avgNumPoints := numPoints / numSegments
		fmt.Printf("Num Tracks: %v; Num Segments: %v (average per track: %v); Num Points: %v (average per segment: %v)\n", len(gpxFile.Tracks), numSegments, avgNumSegments, numPoints, avgNumPoints)
		return []gpx.GPX{}, nil
	}
}
