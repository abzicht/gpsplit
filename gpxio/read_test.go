package gpxio

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileSystem(t *testing.T) {
	gpxFiles, err := ReadFileSystem("../testing/gpxio/gpxio_test_1.gpx")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gpxFiles))
	gpxFile := gpxFiles[0]
	assert.Equal(t, 1, len(gpxFile.Tracks))
	assert.Equal(t, 1, len(gpxFile.Tracks[0].Segments))
	assert.Equal(t, 2, len(gpxFile.Tracks[0].Segments[0].Points))

	gpxFiles, err = ReadFileSystem("../testing/gpxio")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles))
}

func TestRead(t *testing.T) {
	// two concatenated files
	stringReader := strings.NewReader(`
<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
<gpx version="1.1">
<trk>
 <trkseg>
  <trkpt lat="1.14036821" lon="1.54329670">
   <time>1971-01-10T11:00:00Z</time>
   <ele>605.00</ele>
   <sat>22</sat>
   <extensions>
     <speed>0.15</speed>
     <course>300.010</course>
     <accuracy>10.00</accuracy>
     <batterylevel>99.00</batterylevel>
     <useragent>Test Agent</useragent>
   </extensions>
  </trkpt>
 </trkseg>
</trk>
</gpx>
<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
<gpx version="1.1">
<trk>
 <trkseg>
  <trkpt lat="1.24036821" lon="1.84329670">
   <time>1971-01-11T11:00:00Z</time>
   <ele>600.00</ele>
   <sat>32</sat>
   <extensions>
     <speed>0.18</speed>
     <course>270.000</course>
     <accuracy>5.00</accuracy>
     <batterylevel>92.00</batterylevel>
     <useragent>Test Agent</useragent>
   </extensions>
  </trkpt>
 </trkseg>
</trk>
</gpx>
`)
	gpxFiles, err := Read(stringReader)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(gpxFiles))
	gpxFile := gpxFiles[0]
	assert.Equal(t, 1, len(gpxFile.Tracks))
	assert.Equal(t, 1, len(gpxFile.Tracks[0].Segments))
	assert.Equal(t, 1, len(gpxFile.Tracks[0].Segments[0].Points))
}

func TestDelimReader(t *testing.T) {
	stringReader := strings.NewReader(`
<gpx>1</gpx>
<gpx>2</gpx><gpx>3</gpx>
<gpx>4</gpx>
<gpx>4</gp>


<gpx>5</gpx>
`)

	expected := []string{`
<gpx>1</gpx>`, `
<gpx>2</gpx>`,
		`<gpx>3</gpx>`, `
<gpx>4</gpx>`, `
<gpx>4</gp>


<gpx>5</gpx>`,
	}
	r := NewGPXReader(bufio.NewReader(stringReader))

	var err error
	index := 0
	for err == nil {
		var data []byte
		data, err = r.ReadToNextDelim()
		if len(data) == 0 || len(bytes.TrimRight(data, "\n \t")) == 0 {
			continue
		}
		assert.Equal(t, []byte(expected[index]), data)
		index++

	}
	assert.Equal(t, err, io.EOF)
}
