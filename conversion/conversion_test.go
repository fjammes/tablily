package conversion

import (
	"fmt"
	"testing"
)

func TestConvertToLilypond(t *testing.T) {

	tests := []struct {
		fret      int
		stringNum int
		tuning    []string
		duration  int
		deadNote  bool
		rest      bool
		expected  string
	}{
		{3, 1, GuitarTuning, 4, false, false, "g'4\\1"},
		{2, 2, GuitarTuning, 8, false, false, "cis'8\\2"},
		{0, 3, GuitarTuning, 4, false, false, "g4\\3"},
		{1, 4, GuitarTuning, second, false, false, "dis2\\4"},
		{3, 5, GuitarTuning, 4, false, false, "c4\\5"},
		{2, 6, GuitarTuning, 4, false, false, "fis4\\6"},
		{5, 1, BassTuning, 4, false, false, "c'4\\1"},
		{4, 2, BassTuning, 8, false, false, "fis8\\2"},
		{2, 3, BassTuning, 4, false, false, "b4\\3"},
		{3, 4, BassTuning, 2, false, false, "g2\\4"},
		{3, 5, BassTuning, 4, false, true, "r4"},
	}

	for _, test := range tests {
		result, _, err := ConvertToLilypond(test.fret, test.stringNum, test.tuning, test.duration, test.deadNote, test.rest, 0)
		if err != nil {
			t.Errorf("convertToLilypond(%d, %d, %v, %d) = %s; want %s", test.fret, test.stringNum, test.tuning, test.duration, err, test.expected)
		}
		if result != test.expected {
			t.Errorf("convertToLilypond(%d, %d, %v, %d) = %s; want %s", test.fret, test.stringNum, test.tuning, test.duration, result, test.expected)
		}
	}
}

func TestGetNoteIndex(t *testing.T) {
	tests := []struct {
		note     string
		expected int
	}{
		{"C", 0},
		{"C#", 1},
		{"D", 2},
		{"D#", 3},
		{"E", 4},
		{"F", 5},
		{"F#", 6},
		{"G", 7},
		{"G#", 8},
		{"A", 9},
		{"A#", 10},
		{"B", 11},
		{"X", -1},
	}

	for _, test := range tests {
		result, err := getNoteIndex(test.note)
		if err != nil {
			t.Errorf("getNoteIndex(%s) = %s; want %d", test.note, err, test.expected)
		}
		if result != test.expected {
			t.Errorf("getNoteIndex(%s) = %d; want %d", test.note, result, test.expected)
		}
	}
}
func TestParseTabEntry(t *testing.T) {
	tests := []struct {
		tabEntry          string
		previousDuration  int
		previousStringNum int
		expectedFret      int
		expectedDuration  int
		expectedStringNum int
		expectIsDeadNote  bool
		expectRest        bool
		expectError       error
	}{
		{"3:4\\1", sixteenth, 3, 3, 4, 1, false, false, nil},
		{"2\\2", sixteenth, 3, 2, sixteenth, 2, false, false, nil},
		{"0:4", sixteenth, 1, 0, 4, 1, false, false, nil},
		{"12", sixteenth, 1, 12, sixteenth, 1, false, false, nil},
		{"x\\4", sixteenth, 1, 0, sixteenth, 4, true, false, nil},
		// {"3\\5", 0, 16, 3, 5, "", false, false},
		// {"2\\6", 0, 16, 2, 6, "", false, false},
		// {"5\\1", 0, 16, 5, 1, "", false, false},
		// {"4\\2", 0, 16, 4, 2, "", false, false},
		// {"2\\3", 0, 16, 2, 3, "", false, false},
		// {"3\\4", 0, 16, 3, 4, "", false, false},
		// {"3", 4, 16, 3, 4, 16, false, false},
		// {"invalid\\1", 0, "", 0, 0, "", true, false},
		{"3\\invalid", 4, 1, -1, -1, -1, false, false, fmt.Errorf("invalid string number: invalid in note 3\\invalid")},
		// {"", 0, 0, "", 0, "", true, false},
		// {"", 0, 0, "", 0, "", true},
	}

	for _, test := range tests {
		fret, stringNum, duration, isDeadNote, isRest, err := ParseTabEntry(test.tabEntry, test.previousStringNum, test.previousDuration, GuitarTuning)
		fmt.Printf("fret: %d, stringNum: %d, duration: %d, err: %v\n", fret, stringNum, duration, err)
		if test.expectError != nil {
			if err.Error() != test.expectError.Error() {
				t.Errorf("ParseTabEntry(%s, %d, %d) error = '%s', expectError = '%s'", test.tabEntry, test.previousStringNum, test.previousDuration, err.Error(), test.expectError.Error())
			}
		} else if err != nil {
			t.Error("Unexpected error: ", err)
		}
		if fret != test.expectedFret {
			t.Errorf("ParseTabEntry(%s, %d, %d) fret = %d, want %d", test.tabEntry, test.previousStringNum, test.previousDuration, fret, test.expectedFret)
		}
		if duration != test.expectedDuration {
			t.Errorf("ParseTabEntry(%s, %d, %d) duration = %d, want %d", test.tabEntry, test.previousStringNum, test.previousDuration, duration, test.expectedDuration)
		}
		if stringNum != test.expectedStringNum {
			t.Errorf("ParseTabEntry(%s, %d, %d) stringNum = %d, want %d", test.tabEntry, test.previousStringNum, test.previousDuration, stringNum, test.expectedStringNum)
		}
		if isDeadNote != test.expectIsDeadNote {
			t.Errorf("ParseTabEntry(%s, %d, %d) isDeadNote = %v, want %v", test.tabEntry, test.previousStringNum, test.previousDuration, isDeadNote, test.expectIsDeadNote)
		}
		if isRest != test.expectRest {
			t.Errorf("ParseTabEntry(%s, %d, %d) isRest = %v, want %v", test.tabEntry, test.previousStringNum, test.previousDuration, isRest, test.expectRest)
		}
	}
}
