package drum

import (
	"fmt"
	"io"
)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	version string
	tempo   float32
	tracks  []track
}

const stepsPerTrack = 16

// track is the representation of a single track in a Pattern.
type track struct {
	id    uint8
	name  string
	steps [stepsPerTrack]bool
}

func (p *Pattern) read(reader io.Reader) error {
	return nil
}

// String implements the fmt.Stringer interface, allowing Patterns to be printed.
func (p Pattern) String() string {
	// write the header
	str := fmt.Sprintf("Saved with HW Version: %s\nTempo: %g\n", p.version, p.tempo)
	// write each track
	for _, track := range p.tracks {
		str += fmt.Sprintln(track)
	}

	return str
}

// String implements the fmt.Stringer interface, allowing tracks to be printed.
func (t track) String() string {
	// write the header
	header := fmt.Sprintf("(%d) %s\t", t.id, t.name)
	// write the steps
	steps := []byte("|----|----|----|----|")
	// add an 'x' for each note
	for i, x := range t.steps {
		if x {
			// need to adjust 'i' to account for the '|'s
			steps[i+i/4+1] = 'x'
		}
	}
	return header + string(steps)
}
