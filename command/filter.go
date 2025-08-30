package command

import (
	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"gonum.org/v1/gonum/unit"
)

type FilterCommand struct {
	Trim      unit.Length `short:"t" long:"trim" description:"Trim all points at a segment's start AND end that are farther away than the provided radius. If this command is used, --trim-start and --trim-end are ignored." default:"0"`
	TrimStart unit.Length `long:"trim-start" description:"Trim all points at a segment's start that are farther away than the provided radius." default:"0"`
	TrimEnd   unit.Length `long:"trim-end" description:"Trim all points at a segment's end that are farther away than the provided radius." default:"0"`
}

func (f FilterCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	filterOptions := []options.FilterOptions{}
	if f.Trim != 0*unit.Metre {
		filterOptions = append(filterOptions, options.TrimStart(f.Trim))
		filterOptions = append(filterOptions, options.TrimEnd(f.Trim))
	} else {
		if f.TrimStart != 0*unit.Metre {
			filterOptions = append(filterOptions, options.TrimStart(f.TrimStart))
		}
		if f.TrimEnd != 0*unit.Metre {
			filterOptions = append(filterOptions, options.TrimEnd(f.TrimEnd))
		}
	}
	tc = config.NewTransformConfig(config.WithSegmentTransform(gpxtransform.Filter(filterOptions...)))
	return
}
