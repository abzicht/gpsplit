package gpxtransform

import (
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/tkrajina/gpxgo/gpx"
)

/*
SplitFileByTrack creates a "GPXFileTransform"er that splits a file into multiple
files, one for every track contained in the original file.
*/
func SplitFileByTrack() config.GPXFileTransform {
	return func(gpxFile gpx.GPX) (files []gpx.GPX, err error) {
		for trackIndex, _ := range gpxFile.Tracks {
			splitFile := gpxFile
			splitFile.Tracks = []gpx.GPXTrack{gpxFile.Tracks[trackIndex]}
			files = append(files, splitFile)
		}
		return
	}
}

/*
SplitTrackBySegment creates a "GPXTrackTransform"er that splits a track into multiple
tracks, one for every segment contained in the original track.
*/
func SplitTrackBySegment() config.GPXTrackTransform {
	return func(gpxTrack gpx.GPXTrack) (tracks []gpx.GPXTrack, err error) {
		for segmentIndex, _ := range gpxTrack.Segments {
			splitTrack := gpxTrack
			splitTrack.Segments = []gpx.GPXTrackSegment{gpxTrack.Segments[segmentIndex]}
			tracks = append(tracks, splitTrack)
		}
		return
	}
}

/*
Split creates a "GPXSegmentTransform"er that splits a segment into multiple
segments based on splitOptions.
*/
func Split(splitOptions ...options.SplitOptions) config.GPXSegmentTransform {
	return func(trackSegment gpx.GPXTrackSegment) (segments []gpx.GPXTrackSegment, err error) {
		sI, err := findSplitIndices(trackSegment, splitOptions...)
		if err != nil {
			return
		}
		segmentSuffix := &trackSegment
		substractFromIndex := 0
		for _, splitIndex := range sI {
			var segmentPrefix *gpx.GPXTrackSegment
			segmentPrefix, segmentSuffix = segmentSuffix.Split(splitIndex - substractFromIndex)
			substractFromIndex += len(segmentPrefix.Points)
			segments = append(segments, *segmentPrefix)
		}
		segments = append(segments, *segmentSuffix)
		return
	}
}

type splitIndices []int

/*
Returns the indices at that a trackSegment should be split based on the provided options.
*/
func findSplitIndices(trackSegment gpx.GPXTrackSegment, splitOptions ...options.SplitOptions) (sI splitIndices, err error) {
	sI = splitIndices{}
	points := trackSegment.Points
	for index := 0; index < len(points)-1; index++ {
		var doSplit bool = false
		for _, splitOption := range splitOptions {
			doSplit, err = splitOption.Do(trackSegment, index)
			if err != nil {
				return
			}
			if doSplit {
				break
			}
		}
		if doSplit {
			sI = append(sI, index)
		}
	}
	return
}
