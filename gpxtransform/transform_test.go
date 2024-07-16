package gpxtransform

//func TestTransform(t *testing.T) {
//	fileName := "../server/html/static/media/testing.gpx"
//	reader, err := os.Open(fileName)
//	defer reader.Close()
//	assert.NoError(t, err)
//
//	gpxFile, err := gpx.Parse(reader)
//	assert.NoError(t, err)
//	//tc := config.NewTransformConfig(config.WithSegmentTransform(Split(options.TimeSplit(8 * time.Hour))))
//	tc := config.NewTransformConfig(config.WithFileTransform(SplitFileByTrack()), config.WithTrackTransform(SplitTrackBySegment()), config.WithSegmentTransform(Split(options.TimeSplit(8*time.Hour))))
//	gpxFiles, err := TransformFile(*gpxFile, tc)
//	assert.NoError(t, err)
//	fmt.Println(len(gpxFiles))
//
//}
