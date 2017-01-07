package drum

import "testing"
import "fmt"

const expected = `Saved with HW Version: 1.1
Tempo: 120
(0) kick	|x-x-|----|----|----|
(1) snare	|-x-x|----|----|----|
`

func TestPatternString(t *testing.T) {
	p := &Pattern{
		tempo:   120,
		version: "1.1",
		tracks: []track{
			track{
				id:    0,
				name:  "kick",
				steps: [stepsPerTrack]bool{true, false, true, false},
			},
			track{
				id:    1,
				name:  "snare",
				steps: [stepsPerTrack]bool{false, true, false, true},
			},
		},
	}

	actual := fmt.Sprint(p)

	if actual != expected {
		t.Fail()
	}
}
