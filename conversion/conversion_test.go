package conversion

import "testing"

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
