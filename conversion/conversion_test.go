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
		duration  string
		expected  string
	}{
		{3, 1, GuitarTuning, "4", "g'04\\1"},
		{2, 2, GuitarTuning, "8", "b'08\\2"},
		{0, 3, GuitarTuning, "4", "d'04\\3"},
		{1, 4, GuitarTuning, "2", "g'02\\4"},
		{3, 5, GuitarTuning, "4", "d''04\\5"},
		{2, 6, GuitarTuning, "4", "f''04\\6"},
		{5, 1, BassTuning, "4", "a'04\\1"},
		{4, 2, BassTuning, "8", "d'08\\2"},
		{2, 3, BassTuning, "4", "e'04\\3"},
		{3, 4, BassTuning, "2", "a'02\\4"},
	}

	for _, test := range tests {
		result, err := ConvertToLilypond(test.fret, test.stringNum, test.tuning, test.duration)
		if err != nil {
			t.Errorf("convertToLilypond(%d, %d, %v, %s) = %s; want %s", test.fret, test.stringNum, test.tuning, test.duration, err, test.expected)
		}
		if result != test.expected {
			t.Errorf("convertToLilypond(%d, %d, %v, %s) = %s; want %s", test.fret, test.stringNum, test.tuning, test.duration, result, test.expected)
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
		expectError       error
	}{
		{"3:4\\1", 16, 3, 3, 4, 1, nil},
		{"2\\2", 16, 3, 2, 16, 2, nil},
		{"0:4", 16, 1, 0, 4, 1, nil},
		{"12", 16, 1, 12, 16, 1, nil},
		// TODO {"x\\4", 0, 16, 1, 4, -1, nil},
		// {"3\\5", 0, 16, 3, 5, "", false},
		// {"2\\6", 0, 16, 2, 6, "", false},
		// {"5\\1", 0, 16, 5, 1, "", false},
		// {"4\\2", 0, 16, 4, 2, "", false},
		// {"2\\3", 0, 16, 2, 3, "", false},
		// {"3\\4", 0, 16, 3, 4, "", false},
		// {"3", 4, 16, 3, 4, 16, false},
		// {"invalid\\1", 0, "", 0, 0, "", true},
		{"3\\invalid", 4, 1, -1, -1, -1, fmt.Errorf("invalid string number: invalid")},
		// {"", 0, 0, "", 0, "", true},
	}

	for _, test := range tests {
		fret, stringNum, duration, err := ParseTabEntry(test.tabEntry, test.previousStringNum, test.previousDuration)
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
		if stringNum != test.expectedStringNum {
			t.Errorf("ParseTabEntry(%s, %d, %d) stringNum = %d, want %d", test.tabEntry, test.previousStringNum, test.previousDuration, stringNum, test.expectedStringNum)
		}
		if duration != test.expectedDuration {
			t.Errorf("ParseTabEntry(%s, %d, %d) duration = %d, want %d", test.tabEntry, test.previousStringNum, test.previousDuration, duration, test.expectedDuration)
		}
	}
}
