package main

import (
	"log/slog"
	"os"

	"github.com/abzicht/gpsplit/command"
	"github.com/abzicht/gpsplit/gpxio"
	"github.com/abzicht/gpsplit/gpxtransform"
	"github.com/jessevdk/go-flags"
	"github.com/tkrajina/gpxgo/gpx"
)

/*
Reads GPX file(s) and edits them according to the provided command line arguments
*/
func main() {
	var flagOpts command.Flags
	parser := flags.NewParser(&flagOpts, flags.Default)
	parser.LongDescription = "GPSplit is a toolkit for GPX files. \n" +
		"Use the -h flag for more information on sub-commands. Example: `gpsplit filter -h`."
	_, err := parser.Parse()
	if err != nil {
		return
	}

	err = flagOpts.Init()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	tc, err := flagOpts.GetConfiguration(parser.Command.Active.Name)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var gpxFiles []gpx.GPX

	if len(flagOpts.In) == 0 {
		gpxFiles, err = gpxio.Read(os.Stdin)
	} else {
		gpxFiles, err = gpxio.ReadFileSystem(flagOpts.In)
	}
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if len(gpxFiles) == 0 {
		slog.Error("could not read file")
		return
	}

	outFiles, err := gpxtransform.TransformFiles(gpxFiles, tc)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(flagOpts.Out) == 0 {
		err = gpxio.WriteStdout(outFiles)
	} else {
		err = gpxio.WriteFiles(flagOpts.Out, outFiles)
	}
	if err != nil {
		slog.Error(err.Error())
		return
	}

	return
}
