package gpxtransform

import (
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/tkrajina/gpxgo/gpx"
)

func MergeTrackSegments() config.GPXTrackTransform {
	return func(gpxTrack gpx.GPXTrack) (tracks []gpx.GPXTrack, err error) {
		for len(gpxTrack.Segments) > 1 {
			(&gpxTrack).Join(0, 1)
		}
		return []gpx.GPXTrack{gpxTrack}, nil
	}

	// return func(gpxTrack gpx.GPXTrack) (tracks []gpx.GPXTrack, err error) {
	// 	if len(gpxTrack.Segments) == 0 {
	// 		return []gpx.GPXTrack{gpxTrack}, nil
	// 	}

	// 	segment := gpx.GPXTrackSegment{Points: []gpx.GPXPoint{}, Extensions: gpxTrack.Segments[0].Extensions}
	// 	for segmentIndex, _ := range gpxTrack.Segments {
	// 		segment.Points = append(segment.Points, gpxTrack.Segments[segmentIndex].Points...)
	// 	}
	// 	gpxTrack.Segments = []gpx.GPXTrackSegment{segment}
	// 	return []gpx.GPXTrack{gpxTrack}, nil
	// }
}

func MergeTracks() config.GPXFileTransform {
	return func(gpxFile gpx.GPX) (files []gpx.GPX, err error) {
		if len(gpxFile.Tracks) == 0 {
			return []gpx.GPX{gpxFile}, nil
		}
		track := gpxFile.Tracks[0]
		for trackIndex, _ := range gpxFile.Tracks {
			if trackIndex == 0 {
				continue
			}
			track.Segments = append(track.Segments, gpxFile.Tracks[trackIndex].Segments...)
		}
		gpxFile.Tracks = []gpx.GPXTrack{track}
		return []gpx.GPX{gpxFile}, nil
	}
}

func MergeFiles() config.GPXFilesTransform {
	return func(gpxFiles []gpx.GPX) (files []gpx.GPX, err error) {
		if len(gpxFiles) == 0 {
			return
		}
		gpxFile := gpxFiles[0]
		for fileIndex, _ := range gpxFiles {
			if fileIndex == 0 {
				continue
			}
			gpxFile.Tracks = append(gpxFile.Tracks, gpxFiles[fileIndex].Tracks...)
		}
		return []gpx.GPX{gpxFile}, nil
	}
}
