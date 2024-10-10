package command

import (
	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
)

type MergeCommand struct {
	MergeSegments bool `short:"s" long:"merge-segments" description:"Merge multiple segments from a track to a single segment."`
	MergeTracks   bool `short:"t" long:"merge-tracks" description:"Merge multiple tracks to a single track."`
	MergeFiles    bool `short:"f" long:"merge-files" description:"Merge multiple files to a single files."`
}

func (m MergeCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	tc = config.NewTransformConfig()
	if m.MergeSegments {
		tc = config.WithTrackTransform(gpxtransform.MergeTrackSegments())(tc)
	}
	if m.MergeTracks {
		tc = config.WithFileTransform(gpxtransform.MergeTracks())(tc)
	}
	if m.MergeFiles {
		tc = config.WithFilesTransform(gpxtransform.MergeFiles())(tc)
	}
	return
}
