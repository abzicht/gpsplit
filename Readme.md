<img src="gpsplit.svg" alt="GPSplit Logo" style="display: block;margin-left: auto;margin-right: auto;width: 50%;" width="500"/>

# GPSplit

> No more messy GPX files!

GPSplit is a GPX toolkit that is built on top of the
[GPXGO](https://github.com/tkrajina/gpxgo) library by
[@tkrajina](https://github.com/tkrajina).

GPSplit was motivated by a personal GPX file that spans
multiple years of bike tours within a single GPX segment. This GPX file
demanded a fast way of automatically splitting segments, e.g., if there is no
relevant movement over longer time, if single segments are too short, if the
start and end of a segment has no significant movement, or if there is a huge
time / distance jump between two points.

* GPSplit provides functionality for splitting, filtering, and analyzing GPX
  files, tracks, and segments.
* The GPSplit command line interface allows combining multiple split / filter
  rules through UNIX pipes.
* GPSplit makes it easy to define custom split / filter rules.
* GPSplit can be used as a library in combination with
  [GPXGO](https://github.com/tkrajina/gpxgo).

## Installation
Install Go 1.22 and run
```bash
go install github.com/abzicht/gpsplit
```

## Examples

This command splits a GPX file into multiple files. Splitting is performed, if
a GPX segment has a time jump of over 8 hours between two points.

```bash
cat ./my-recording.gpx | gpsplit --out ./ split --time 8h
```

Now, what if we want to detect pauses, split the segments at those, then remove
those pauses, thereafter simplify all polylines, and finally split every
segment into its own file? Then this is the command to use:

```bash
gpsplit --in ./my-recording.gpx split --pause-split 200,1h | gpsplit filter --trim 50 | gpsplit direct --simplify 40 | gpsplit --out ./ split
```

While `split --pause-split 200,1h` splits into multiple segments whenever we
are in a 200 meter radius for over one hour, `filter --trim 50` removes all
points at the start and end of a segment that are in a 50 meter radius. `direct
--simplify 40` then applies the Ramer-Douglas-Peucker algorithm for removing
unnecessary GPX points. Finally, the singular, multi-segmented GPX file is
split into individual files per segment that are stored in `./` with `--out ./
split`.

## Results
The gif below shows the before (two GPX segments, spanning multiple years) and after (over 100 individual tracks) of using `gpsplit`.

![Example showing before and after treatment with gpsplit](before-and-after.gif)
