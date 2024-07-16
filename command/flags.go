package command

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/abzicht/gpsplit/gpxtransform/config"
)

/*
Flag holds all arguments passed via command line
*/
type Flags struct {
	In      string         `short:"i" long:"in" description:"The single file that new GPX data is read from. Leave empty to read from STDIN"`
	Out     string         `short:"o" long:"out" description:"The file or folder that new GPX data is written to. Leave empty to write to STDOUT"`
	Verbose []bool         `short:"v" long:"verbosity" description:"Verbosity with that information is printed to STDOUT"`
	Split   SplitCommand   `command:"split" description:"Splits GPX segments into multiple segments or files"`
	Filter  FilterCommand  `command:"filter" description:"Applies filters on GPX points"`
	Direct  DirectCommand  `command:"direct" description:"Applies misc. functionality directly on GPX segments"`
	Analyze AnalyzeCommand `command:"analyze" description:"Prints information for the provided GPX data"`
}

/*
Init performs initialization, including setting a global slog-level
*/
func (flagOpts Flags) Init() (err error) {
	// See: https://pkg.go.dev/log/slog#Level
	level := slog.LevelError
	switch max(-4, min(8, 8-4*len(flagOpts.Verbose))) {
	case int(slog.LevelDebug):
		level = slog.LevelDebug
	case int(slog.LevelInfo):
		level = slog.LevelInfo
	case int(slog.LevelWarn):
		level = slog.LevelWarn
	case int(slog.LevelError):
		level = slog.LevelError
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

/*
GetConfiguration returns the configuration for the sub-command specified with name.
Valid command names are split, filter, direct, and analyze.
*/
func (flagOpts Flags) GetConfiguration(name string) (tc config.TransformConfig, err error) {
	switch name {
	case "split":
		tc, err = flagOpts.Split.GetConfiguration()
	case "filter":
		tc, err = flagOpts.Filter.GetConfiguration()
	case "direct":
		tc, err = flagOpts.Direct.GetConfiguration()
	case "analyze":
		tc, err = flagOpts.Analyze.GetConfiguration()
	default:
		err = errors.New(fmt.Sprintf("unknown command: %v", name))
		return
	}
	return
}
