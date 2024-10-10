package gpxtransform

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/tkrajina/gpxgo/gpx"
)

type TransformError struct {
	Msg string
}

func (te TransformError) Error() string {
	return te.Msg
}

/* TransformSegment applies a provided config.TransformConfig on a gpx track
* segment. It returns zero, one, or multiple track segments depending on the
* applied transformation.
 */
func TransformSegment(segment gpx.GPXTrackSegment, tc config.TransformConfig) ([]gpx.GPXTrackSegment, error) {
	switch tc.SegmentT.(type) {
	case fun.Some[config.GPXSegmentTransform]:
		// action can be a filter (returning 0 or 1 elements), a map (modifying the
		// segment), or a splitter (returning 1 or more elements)
		return tc.SegmentT.GetValue()(segment)
	default:
		return []gpx.GPXTrackSegment{segment}, nil
	}
}

/*
TransformTrack applies a provided config.TransformConfig on a gpx track. It
* returns zero, one, or multiple tracks depending on the
* applied transformation.
*/
func TransformTrack(track gpx.GPXTrack, tc config.TransformConfig) ([]gpx.GPXTrack, error) {
	segments := []gpx.GPXTrackSegment{}
	for segmentIndex, _ := range track.Segments {
		s, err := TransformSegment(track.Segments[segmentIndex], tc)
		if err != nil {
			return nil, err
		}
		segments = append(segments, s...)
	}
	track.Segments = segments
	switch tc.TrackT.(type) {
	case fun.Some[config.GPXTrackTransform]:
		return tc.TrackT.GetValue()(track)
	default:
		return []gpx.GPXTrack{track}, nil
	}
}

/*
TransformFile applies a provided config.TransformConfig on a gpx file. It
* returns zero, one, or multiple files depending on the
* applied transformation.
*/
func TransformFile(gpxFile gpx.GPX, tc config.TransformConfig) ([]gpx.GPX, error) {
	tracks := []gpx.GPXTrack{}
	for trackIndex, _ := range gpxFile.Tracks {
		t, err := TransformTrack(gpxFile.Tracks[trackIndex], tc)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t...)
	}
	gpxFile.Tracks = tracks
	switch tc.FileT.(type) {
	case fun.Some[config.GPXFileTransform]:
		return tc.FileT.GetValue()(gpxFile)
	default:
		return []gpx.GPX{gpxFile}, nil
	}
}

/*
TransformFiles applies a provided config.TransformConfig on multiple gpx files. It
* returns zero, one, or multiple files depending on the
* applied transformation.
*/
func TransformFiles(gpxFiles []gpx.GPX, tc config.TransformConfig) ([]gpx.GPX, error) {
	files := []gpx.GPX{}
	for fileIndex, _ := range gpxFiles {
		f, err := TransformFile(gpxFiles[fileIndex], tc)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	gpxFiles = files
	switch tc.FilesT.(type) {
	case fun.Some[config.GPXFilesTransform]:
		return tc.FilesT.GetValue()(gpxFiles)
	default:
		return gpxFiles, nil
	}
}
