package gpxio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/tkrajina/gpxgo/gpx"
)

/* Write all gpx objects to STDOUT.
 * Files are pretty-printed, i.e., use indentation.
 * Files are separated by a single empty line.
 */
func WriteStdout(outFiles []gpx.GPX) (err error) {
	for _, gpxFile := range outFiles {
		var xmlBytes []byte
		xmlBytes, err = gpxFile.ToXml(gpx.ToXmlParams{Version: gpxFile.Version, Indent: true})
		if err != nil {
			return
		}
		var n int
		n, err = os.Stdout.Write(xmlBytes)
		if err != nil {
			return
		}
		if n != len(xmlBytes) {
			err = errors.New("failed to write to STDOUT")
			return
		}
		n, err = os.Stdout.Write([]byte("\n"))
		if err != nil {
			return
		}
		if n != 1 {
			err = errors.New("failed to write to STDOUT")
			return
		}
	}
	return
}

/* Save all gpx objects to files.
 * If the given fileName points to a folder: all files are named according to
 * the name defined in the gpx object.
 * Else: all files use the fileName as prefix.
 */
func WriteFiles(fileName string, outFiles []gpx.GPX) (err error) {

	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist. That error is expected and does not need to
			// be escalated.
			info = nil
			err = nil
		} else {
			return
		}
	}

	var getFileName func(gpxFile gpx.GPX, index int) string
	if nil != info && info.IsDir() {
		getFileName = func(gpxFile gpx.GPX, index int) string {
			baseName := slug.Make(gpxFile.Name)
			return filepath.Join(fileName, fmt.Sprintf("%v-%v.gpx", baseName, index))
		}
	} else {
		getFileName = func(gpxFile gpx.GPX, index int) string {
			return fmt.Sprintf("%v-%v.gpx", fileName, index)
		}
	}

	for i, gpxFile := range outFiles {
		xmlBytes, err2 := gpxFile.ToXml(gpx.ToXmlParams{Version: gpxFile.Version, Indent: true})
		if err2 != nil {
			err = err2
			return
		}
		writer, err2 := os.Create(getFileName(gpxFile, i+1))
		if err2 != nil {
			err = err2
			return
		}
		n, err2 := writer.Write(xmlBytes)
		if err2 != nil {
			err = err2
			return
		}
		if n != len(xmlBytes) {
			err = errors.New(fmt.Sprintf("could not store all data in %v", getFileName(gpxFile, i+1)))
			return
		}
	}
	return
}
