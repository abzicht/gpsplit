package options

import (
	"github.com/tkrajina/gpxgo/gpx"
)

/*
TransformOptions is passed to transformation functions.
It holds a function that decides whether a given transformation should be performed.
*/
type TransformOptions struct {
	/*
		do must return true whenever a transformation should be performed at the provided index of a given track segment.
	*/
	Do func(segment gpx.GPXTrackSegment, index int) (bool, error)
}

/*
NewTransformOptions returns a new TransformOptions object that holds the provided do function.
*/
func NewTransformOptions(do func(gpx.GPXTrackSegment, int) (bool, error)) TransformOptions {
	return TransformOptions{do}
}
