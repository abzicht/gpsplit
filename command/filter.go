package command

import (
	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"gonum.org/v1/gonum/unit"
)

type FilterCommand struct {
	Trim      unit.Length `short:"t" long:"trim" description:"Trim all points at a segment's start AND end that are farther away than the provided radius. If this command is used, --trim-left and --trim-right are ignored." default:"0"`
	TrimLeft  unit.Length `long:"trim-left" description:"Trim all points at a segment's start that are farther away than the provided radius." default:"0"`
	TrimRight unit.Length `long:"trim-right" description:"Trim all points at a segment's end that are farther away than the provided radius." default:"0"`
}

func (f FilterCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	filterOptions := []options.FilterOptions{}
	if f.Trim != 0*unit.Metre {
		filterOptions = append(filterOptions, options.TrimLeft(f.Trim))
		filterOptions = append(filterOptions, options.TrimRight(f.Trim))
	} else {
		if f.TrimLeft != 0*unit.Metre {
			filterOptions = append(filterOptions, options.TrimLeft(f.TrimLeft))
		}
		if f.TrimRight != 0*unit.Metre {
			filterOptions = append(filterOptions, options.TrimRight(f.TrimRight))
		}
	}
	tc = config.NewTransformConfig(config.WithSegmentTransform(gpxtransform.Filter(filterOptions...)))
	return
}
