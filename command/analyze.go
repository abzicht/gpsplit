package command

import (
	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
)

type AnalyzeCommand struct {
	Num bool `short:"c" long:"count" description:"Only analyze object counts."`
}

func (a AnalyzeCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	tc = config.NewTransformConfig()
	if a.Num {
		tc = config.WithFileTransform(gpxtransform.CountFile())(tc)
	} else {
		tc = config.WithFileTransform(gpxtransform.AnalyzeFile())(tc)
	}
	return
}
