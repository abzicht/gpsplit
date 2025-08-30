package command

import (
	"time"

	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"gonum.org/v1/gonum/unit"
)

type DirectCommand struct {
	Simplify    unit.Length   `short:"s" long:"simplify" description:"Simplify tracks using Ramer-Douglas-Peucker algorithm with the provided distance in meters." default:"0"`
	MinPoints   int           `long:"min-points" description:"Remove segments that have less points than the provided number." default:"0"`
	MinRadius   unit.Length   `long:"min-radius" description:"Remove segments whose points are all within a given radius from the starting point." default:"0"`
	MinDistance unit.Length   `long:"min-distance" description:"Remove segments that are shorter than the provided min distance." default:"0"`
	MinDuration time.Duration `long:"min-duration" description:"Remove segments that are shorter than the provided min duration." default:"0"`
	MaxDuration time.Duration `long:"max-duration" description:"Remove segments that are longer than the provided max duration." default:"0"`
	RemoveStops bool          `long:"remove-stops" description:"Remove points that are considered to have no movement. Cf. command 'split' for splitting segments at stops."`
}

func (d DirectCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	defaultDuration, err := time.ParseDuration("0s")
	if err != nil {
		return
	}

	directOptions := []options.DirectOptions{}
	if d.Simplify != 0*unit.Metre {
		directOptions = append(directOptions, options.DouglasPeucker(d.Simplify))
	}
	if d.MinPoints != 0 {
		directOptions = append(directOptions, options.MinPoints(d.MinPoints))
	}
	if d.MinDistance != 0*unit.Metre {
		directOptions = append(directOptions, options.MinDistance(d.MinDistance))
	}
	if d.MinRadius != 0*unit.Metre {
		directOptions = append(directOptions, options.MinRadius(d.MinRadius))
	}
	if d.MinDuration != defaultDuration {
		directOptions = append(directOptions, options.MinDuration(d.MinDuration))
	}
	if d.MaxDuration != defaultDuration {
		directOptions = append(directOptions, options.MaxDuration(d.MaxDuration))
	}
	if d.RemoveStops {
		directOptions = append(directOptions, options.RemoveStops())
	}
	tc = config.NewTransformConfig(config.WithSegmentTransform(gpxtransform.Direct(directOptions...)))
	return
}
