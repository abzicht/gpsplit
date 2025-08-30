package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/abzicht/gpsplit/gpxtransform/config"
	"github.com/abzicht/gpsplit/gpxtransform/options"
	"gonum.org/v1/gonum/unit"
)

type SplitCommand struct {
	Tracks     bool          `short:"t" long:"tracks" description:"Create an individual file for every track."`
	Segments   bool          `short:"s" long:"segments" description:"Create an individual track for every track segment."`
	Distance   unit.Length   `short:"d" long:"distance" description:"If the distance between two consecutive points is larger than this value (in meters), the path is split between those points." default:"0"`
	Duration   time.Duration `short:"u" long:"duration" description:"If the time between two consecutive points is larger than this duration, the path is split between those points. Valid formats: 5m (5 minutes); 12h (12 hours); 4h12m0s (4 hours 12 minutes 0 seconds)." default:"0s"`
	PauseSplit string        `long:"pause-split" description:"Split paths, if no movement over a longer period is detected, i.e., if points lie within a provided radius for a given minimal time. Time and radius are comma-separated in strict order (RADIUS,TIME). Valid formats: 20,5m (20 meters, 5 minutes); 1000,12h (1 kilometer, 12 hours); 300,4h12m0s (300 meters, 4 hours 12 minutes 0 seconds)." default:""`
}

func (s SplitCommand) GetConfiguration() (tc config.TransformConfig, err error) {
	defaultDuration, err := time.ParseDuration("0s")
	if err != nil {
		return
	}
	splitOptions := []options.SplitOptions{}
	// TimeSplit
	if s.Duration != defaultDuration {
		splitOptions = append(splitOptions, options.TimeSplit(s.Duration))
	}
	// DistanceSplit
	if s.Distance != 0*unit.Metre {
		splitOptions = append(splitOptions, options.DistanceSplit(s.Distance))
	}
	// PauseSplit
	if 0 != len(s.PauseSplit) {
		radiusAndTime := strings.Split(s.PauseSplit, ",")
		if 2 != len(radiusAndTime) {
			err = CommandError{"incorrect format for no-pause; expecting format RADIUS,TIME"}
			return
		}
		var radiusInt int
		radiusInt, err = strconv.Atoi(radiusAndTime[0])
		if err != nil {
			err = CommandError{"incorrect format for radius in no-pause"}
			return
		}
		var duration time.Duration
		duration, err = time.ParseDuration(radiusAndTime[1])
		if err != nil {
			return
		}
		splitOptions = append(splitOptions, options.PauseSplit(unit.Length(radiusInt)*unit.Metre, duration))
	}
	tc = config.NewTransformConfig(config.WithSegmentTransform(gpxtransform.Split(splitOptions...)))
	if s.Tracks {
		tc = config.WithFileTransform(gpxtransform.SplitFileByTrack())(tc)
	}
	if s.Segments {
		tc = config.WithTrackTransform(gpxtransform.SplitTrackBySegment())(tc)
	}
	return
}
