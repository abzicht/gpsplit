package config

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/tkrajina/gpxgo/gpx"
)

type GPXFileTransform func(gpxFile gpx.GPX) ([]gpx.GPX, error)
type GPXFilesTransform func(gpxFiles []gpx.GPX) ([]gpx.GPX, error)
type GPXTrackTransform func(gpxTrack gpx.GPXTrack) ([]gpx.GPXTrack, error)
type GPXSegmentTransform func(gpxSegment gpx.GPXTrackSegment) ([]gpx.GPXTrackSegment, error)
type TransformConfig struct {
	FileT    fun.Option[GPXFileTransform]
	FilesT   fun.Option[GPXFilesTransform]
	TrackT   fun.Option[GPXTrackTransform]
	SegmentT fun.Option[GPXSegmentTransform]
}

type TransformConfigOpt func(tc TransformConfig) TransformConfig

func WithFileTransform(t GPXFileTransform) TransformConfigOpt {
	return func(tc TransformConfig) TransformConfig {
		tc.FileT = fun.NewSome[GPXFileTransform](t)
		return tc
	}
}

func WithFilesTransform(t GPXFilesTransform) TransformConfigOpt {
	return func(tc TransformConfig) TransformConfig {
		tc.FilesT = fun.NewSome[GPXFilesTransform](t)
		return tc
	}
}

func WithTrackTransform(t GPXTrackTransform) TransformConfigOpt {
	return func(tc TransformConfig) TransformConfig {
		tc.TrackT = fun.NewSome[GPXTrackTransform](t)
		return tc
	}
}

func WithSegmentTransform(t GPXSegmentTransform) TransformConfigOpt {
	return func(tc TransformConfig) TransformConfig {
		tc.SegmentT = fun.NewSome[GPXSegmentTransform](t)
		return tc
	}
}

func NewTransformConfig(opts ...TransformConfigOpt) TransformConfig {

	tc := TransformConfig{
		FileT:    fun.NewNone[GPXFileTransform](),
		FilesT:   fun.NewNone[GPXFilesTransform](),
		TrackT:   fun.NewNone[GPXTrackTransform](),
		SegmentT: fun.NewNone[GPXSegmentTransform](),
	}
	for _, opt := range opts {
		tc = opt(tc)
	}
	return tc
}
