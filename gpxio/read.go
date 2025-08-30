package gpxio

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tkrajina/gpxgo/gpx"
)

/*
ReadFileSystem reads file(s) from the file system using the provided path (fileName).
If fileName is a folder, it reads all top-level files ending with ".gpx" or ".GPX" in that folder
*/
func ReadFileSystem(fileName string) (gpxFiles []gpx.GPX, err error) {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		err = errors.Join(err, errors.New(fmt.Sprintf("directory or file does not exist: \"%v\"", fileName)))
		return
	} else if err != nil {
		return
	}
	if info.IsDir() {
		return ReadFolder(fileName)
	} else {
		var gpxFile gpx.GPX
		gpxFile, err = ReadFile(fileName)
		if err != nil {
			return
		}
		gpxFiles = []gpx.GPX{gpxFile}
		return
	}
}

/*
ReadFile reads a single file from the file system that is identified with the provided fileName
*/
func ReadFile(fileName string) (gpxFile gpx.GPX, err error) {
	reader, err := os.Open(fileName)
	if err != nil {
		err = errors.Join(errors.New(fmt.Sprintf("could not open file %v", fileName)), err)
		return
	}
	gpxFilePointer, err := gpx.Parse(reader)
	gpxFile = *gpxFilePointer
	return
}

/*
ReadFolder reads the GPX files (ending with ".gpx" or ".GPX") of a folder (without recursion) and returns their contents
*/
func ReadFolder(folderName string) (gpxFiles []gpx.GPX, err error) {
	dirEntries, err := os.ReadDir(folderName)
	gpxFiles = []gpx.GPX{}
	if err != nil {
		err = errors.Join(errors.New(fmt.Sprintf("could not read folder %v", folderName)), err)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".gpx") {
			continue
		}
		var gpxFile gpx.GPX
		gpxFile, err = ReadFile(filepath.Join(folderName, entry.Name()))
		if err != nil {
			return
		}
		gpxFiles = append(gpxFiles, gpxFile)
	}
	return
}

/*
Read reads gpx data that is written, e.g., to STDIN.
If multiple gpx files are written to the reader, it separates those based on the ending tag "</gpx>"
*/
func Read(r io.Reader) (gpxFiles []gpx.GPX, err error) {
	reader := bufio.NewReader(r)

	sReader := NewGPXReader(reader)
	for {
		var data []byte
		data, err = sReader.ReadToNextDelim()
		if err != nil && err != io.EOF {
			break
		}
		if len(data) == 0 || len(bytes.TrimRight(data, "\n \t")) == 0 {
			break
		}
		gpxFile, err2 := gpx.Parse(bytes.NewBuffer(data))
		if err2 != nil {
			err = err2
			return
		}
		gpxFiles = append(gpxFiles, *gpxFile)
		if err == io.EOF {
			break
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

/*
A custom reader that tries to separate, e.g., multiple GPX files from a single stream.
This is done by searching for a delimiter.
*/
type DelimReader struct {
	reader *bufio.Reader
	buffer []byte
	delim  []byte
}

/*
NewDelimReader creates a new DelimReader that reads from reader and splits
based on the provided delimiter.

The delimiter's last byte must not occur
in previous bytes of the delimiter. This works, e.g., if delim is
"</gpx>", as ">" only occurs once.
*/
func NewDelimReader(reader *bufio.Reader, delimiter []byte) DelimReader {
	return DelimReader{reader, []byte{}, delimiter}
}

/*
NewGPXReader creates a new DelimReader that reads from reader and splits
concatenated GPX files.
*/
func NewGPXReader(reader *bufio.Reader) DelimReader {
	return DelimReader{reader, []byte{}, []byte("</gpx>")}
}

/*
splitWithDelim splits a byte sequence, if it contains the reader's configured delimiter.
all bytes up to and including the delimiter are returned in "head", the rest is returned in "tail".
"split" is true, iff the delimiter was found.
*/
func (s DelimReader) splitWithDelim(p []byte) (head, tail []byte, split bool) {
	delimIndex := bytes.Index(p, s.delim)
	if delimIndex == -1 {
		return p, []byte{}, false
	}
	return p[0 : delimIndex+len(s.delim)], p[delimIndex+len(s.delim) : len(p)], true

}

/*
readDelim reads until the next potential delimiter and return the bytes up to that point.
delimReached indicates, if we actually reached the delimiter.
*/
func (s DelimReader) readDelim() (p []byte, delimReached bool, err error) {
	delimReached = false
	if len(s.buffer) != 0 {
		var head, tail []byte
		head, tail, delimReached = s.splitWithDelim(s.buffer)
		if delimReached {
			p = head
			s.buffer = tail
		} else {
			p = s.buffer
		}
		return
	}
	data, err := s.reader.ReadSlice(s.delim[len(s.delim)-1])
	if err != nil && err != io.EOF {
		return
	}
	head, tail, delimReached := s.splitWithDelim(data)
	if delimReached {
		// we we are at a delimiter. We only return the data up to and including the
		// delimiter. The rest is buffered and returned with the next
		// call.
		p = head
		s.buffer = tail
	} else {
		// we did not arrive at a delimiter, so we continue
		p = data
	}
	return
}

/*
ReadToNextDelim reads to the next full delimiter and returns all bytes that
came before (including the delimiter). Reads until EOF and also returns the
rest between the last delimiter and EOF. Returns err as io.EOF, if the
end of the file is reached.
*/
func (s DelimReader) ReadToNextDelim() (p []byte, err error) {
	buffer := bytes.NewBuffer([]byte{})
	for {
		var data []byte
		var delimReached bool
		data, delimReached, err = s.readDelim()
		if err != nil && err != io.EOF {
			return buffer.Bytes(), err
		}
		_, err2 := buffer.Write(data)
		if err2 != nil {
			return buffer.Bytes(), err2
		}
		if delimReached || err == io.EOF {
			return buffer.Bytes(), err
		}
	}
}
