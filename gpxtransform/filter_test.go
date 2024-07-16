package gpxtransform

import (
	"testing"

	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/stretchr/testify/assert"
	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

func TestFilterTrimLeft(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithSegmentTransform(Filter(options.TrimLeft(10 * unit.Metre))))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 2, len(gpxFiles[0].Tracks[0].Segments[0].Points))

	tc = config.NewTransformConfig(config.WithSegmentTransform(Filter(options.TrimLeft(100 * unit.Metre))))
	gpxFiles, err = TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 1, len(gpxFiles[0].Tracks[0].Segments[0].Points))

}

func TestFilterTrimRight(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithSegmentTransform(Filter(options.TrimRight(10 * unit.Metre))))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 2, len(gpxFiles[0].Tracks[0].Segments[0].Points))

	tc = config.NewTransformConfig(config.WithSegmentTransform(Filter(options.TrimRight(100 * unit.Metre))))
	gpxFiles, err = TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 1, len(gpxFiles[0].Tracks[0].Segments[0].Points))

}
