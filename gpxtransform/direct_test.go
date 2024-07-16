package gpxtransform

import (
	"testing"

	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/stretchr/testify/assert"
	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

func TestDirect(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	// add again, when having gpx data with time
	//tc := config.NewTransformConfig(config.WithSegmentTransform(Direct(options.MinDuration(10 * time.Minute))))
	tc := config.NewTransformConfig(config.WithSegmentTransform(Direct(options.MinDistance(10 * unit.Metre))))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 2, len(gpxFiles[0].Tracks[0].Segments))

	tc = config.NewTransformConfig(config.WithSegmentTransform(Direct(options.MinDistance(100 * unit.Metre))))
	gpxFiles, err = TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 0, len(gpxFiles[0].Tracks[0].Segments))
}
