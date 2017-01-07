package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// Decoder reads and decodes a Pattern from file.
type Decoder struct {
	r    io.Reader
	path string
	// size of the encoded data
	size int64
}

// NewDecoder creates new Pattern Decoder from a file path
func NewDecoder(path string) *Decoder {
	return &Decoder{path: path}
}

// Decode decodes contents of d.path file and returns Pattern object
func (d *Decoder) Decode() (*Pattern, error) {
	f, err := os.Open(d.path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	d.r = f

	pattern := &Pattern{}
	err = d.readHeader(pattern)
	if err != nil {
		return nil, err
	}
	err = d.readTracks(pattern)
	if err != nil {
		return nil, err
	}

	return pattern, err
}

func (d *Decoder) readHeader(p *Pattern) error {
	// SPLICE
	splice := make([]byte, 6)
	err := binary.Read(d.r, binary.BigEndian, &splice)
	if err != nil {
		return err
	}

	if string(splice) != "SPLICE" {
		return errors.New("invalid file format: no SPLICE header")
	}

	// 8 byte int
	var size int64
	err = binary.Read(d.r, binary.BigEndian, &size)
	if err != nil {
		return err
	}
	if size < 36 {
		return errors.New("invalid data size")
	}
	d.size = size - 36 // size of data except version and tempo

	// 32 byte version string
	ver := make([]byte, 32)
	err = binary.Read(d.r, binary.BigEndian, &ver)
	if err != nil {
		return err
	}
	// trim zeroes from version string
	p.version = string(bytes.TrimRight(ver[:], string(0)))

	// tempo is little endian float32
	var tempo float32
	err = binary.Read(d.r, binary.LittleEndian, &tempo)
	if err != nil {
		return err
	}
	p.tempo = tempo

	return nil
}

func (d *Decoder) readTracks(p *Pattern) error {
	for d.size > 0 {
		track := &track{}
		err := binary.Read(d.r, binary.BigEndian, &track.id)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return errors.New("wrong track header: " + err.Error())
		}

		var trackNameLen uint32
		err = binary.Read(d.r, binary.BigEndian, &trackNameLen)
		if err != nil {
			return errors.New("wrong track header: " + err.Error())
		}

		trackName := make([]byte, trackNameLen)
		err = binary.Read(d.r, binary.BigEndian, &trackName)
		if err != nil {
			return errors.New("wrong track header: " + err.Error())
		}
		track.name = string(trackName[:])

		stepsBuffer := make([]byte, stepsPerTrack)
		err = binary.Read(d.r, binary.BigEndian, &stepsBuffer)
		if err != nil {
			return errors.New("wrong track steps: " + err.Error())
		}
		var steps [stepsPerTrack]bool
		for i, v := range stepsBuffer {
			if v > 1 {
				return errors.New("illegal step value")
			}
			steps[i] = v == 1
		}
		track.steps = steps

		p.tracks = append(p.tracks, *track)
		// 21 = 16 byte steps + 4 byte ID + 1 byte nameLen
		d.size -= 21 + int64(trackNameLen)
	}

	return nil
}
