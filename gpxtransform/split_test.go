package gpxtransform

import (
	"testing"

	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"github.com/stretchr/testify/assert"
	"github.com/tkrajina/gpxgo/gpx"
	"gonum.org/v1/gonum/unit"
)

const (
	gpxData = `<?xml version="1.0"?>
<gpx creator="GPS Visualizer http://www.gpsvisualizer.com/" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">
<trk>
  <name>127 A27 - 21 Wharf Rd</name>
  <desc>21.7 km, 1:10</desc>
  <trkseg>
    <trkpt lat="50.87551" lon="-1.28259">
      <ele>46.164</ele>
    </trkpt>
    <trkpt lat="50.87545" lon="-1.28237">
      <ele>46.848</ele>
    </trkpt>
  </trkseg>
  <trkseg>
    <trkpt lat="50.87533" lon="-1.28199">
      <ele>47.812</ele>
    </trkpt>
    <trkpt lat="50.87522" lon="-1.28158">
      <ele>48.0</ele>
    </trkpt>
  </trkseg>
</trk>
<trk>
  <name>128 A28 - 22 Wherf Rd</name>
  <desc>22.7 km, 1:20</desc>
  <trkseg>
    <trkpt lat="50.87533" lon="-1.28199">
      <ele>47.812</ele>
    </trkpt>
    <trkpt lat="50.87522" lon="-1.28158">
      <ele>48.0</ele>
    </trkpt>
  </trkseg>
  <trkseg>
    <trkpt lat="50.87551" lon="-1.28259">
      <ele>46.164</ele>
    </trkpt>
    <trkpt lat="50.87545" lon="-1.28237">
      <ele>46.848</ele>
    </trkpt>
  </trkseg>
</trk>
</gpx>
`

	trackNameA = "127 A27 - 21 Wharf Rd"
	trackNameB = "128 A28 - 22 Wherf Rd"
)

func TestSplitFileByTrack(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithFileTransform(SplitFileByTrack()))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles))

	assert.Equal(t, trackNameA, gpxFiles[0].Tracks[0].Name)
	assert.Equal(t, trackNameB, gpxFiles[1].Tracks[0].Name)
}
func TestSplitTrackBySegment(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithTrackTransform(SplitTrackBySegment()))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gpxFiles))
	assert.Equal(t, 4, len(gpxFiles[0].Tracks))

	assert.Equal(t, trackNameA, gpxFiles[0].Tracks[0].Name)
	assert.Equal(t, trackNameA, gpxFiles[0].Tracks[1].Name)
	assert.Equal(t, trackNameB, gpxFiles[0].Tracks[2].Name)
	assert.Equal(t, trackNameB, gpxFiles[0].Tracks[3].Name)
}

func TestSplitFileAndTrack(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithFileTransform(SplitFileByTrack()), config.WithTrackTransform(SplitTrackBySegment()))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(gpxFiles))
}

func TestSplit(t *testing.T) {
	gpxFile, err := gpx.ParseBytes([]byte(gpxData))
	assert.NoError(t, err)
	tc := config.NewTransformConfig(config.WithSegmentTransform(Split(options.DistanceSplit(1 * unit.Metre))))
	gpxFiles, err := TransformFile(*gpxFile, tc)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gpxFiles))
	assert.Equal(t, 2, len(gpxFiles[0].Tracks))
	assert.Equal(t, 4, len(gpxFiles[0].Tracks[0].Segments))
	assert.Equal(t, 1, len(gpxFiles[0].Tracks[0].Segments[0].Points))
}
