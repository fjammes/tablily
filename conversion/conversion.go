package conversion

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Exported standard tuning for guitar (E A D G B e)
var GuitarTuning = []string{"e,", "a,", "d", "g", "b", "e'"}

// Exported standard tuning for bass (E A D G)
var BassTuning = []string{"e,", "a,", "d", "g"}

var noteMap = map[string]int{
	"e,":  -8,
	"a,":  -3,
	"c":   0,
	"cis": 1,
	"d":   2,
	"dis": 3,
	"e":   4,
	"f":   5,
	"fis": 6,
	"g":   7,
	"gis": 8,
	"a":   9,
	"ais": 10,
	"b":   11,
	"e'":  16,
}

// getOpenStringNote returns the note of the open string for the given string number
// stringNum is the string number (1 is the highest string)
func getOpenStringNote(tuning []string, stringNum int) (string, error) {
	numberOfString := len(tuning)
	if stringNum < 1 || stringNum > numberOfString {
		return "", fmt.Errorf("invalid string number: %d", stringNum)
	}
	return tuning[numberOfString-stringNum], nil
}

func getNoteIndex(note string) (int, error) {
	scaleIndex := -1
	var err error
	// Same as below but with a map

	if val, ok := noteMap[note]; ok {
		scaleIndex = val
	} else {
		err = fmt.Errorf("invalid note: %s", note)
	}
	return scaleIndex, err
}

// ConvertToLilypond converts a tab entry to LilyPond format
func ConvertToLilypond(fret, stringNum int, tuning []string, duration string) (string, error) {

	slog.Debug("Convert tab entry to LilyPond format", "fret", fret, "stringNum", stringNum, "tuning", tuning, "duration", duration)

	// Calculate the note name and octave
	openStringNote, err := getOpenStringNote(tuning, stringNum)
	if err != nil {
		slog.Error("Invalid string number", "stringNum", stringNum)
		return "", err
	}
	openStringNoteIndex, err := getNoteIndex(openStringNote)
	if err != nil {
		slog.Error("Invalid open string note", "openStringNote", openStringNote)
		return "", err
	}
	noteIndex := openStringNoteIndex + fret

	slog.Debug("Calculate the note index", "openStringNote", openStringNote, "openStringNoteIndex", openStringNoteIndex, "noteIndex", noteIndex)

	noteNames, err := findNoteName(noteIndex)
	if err != nil {
		slog.Error("Invalid note index", "noteIndex", noteIndex)
		return "", err
	}

	// Convert duration to an integer to remove leading zeros
	durationInt, _ := strconv.Atoi(duration)
	return fmt.Sprintf("%s%d\\%d", noteNames, durationInt, stringNum), nil
}

func findNoteName(noteIndex int) (string, error) {
	var octave int
	var octaveShift string
	noteName := ""

	if noteIndex < 0 {
		octave = -1
	}

	scaleIndex := noteIndex % 12
	for note, i := range noteMap {
		if i == scaleIndex {
			noteName = note
		}
	}
	if noteName == "" {
		return "", fmt.Errorf("invalid note index: %d", noteIndex)
	}
	octave = Abs(noteIndex) / 12
	if noteIndex < 0 {
		octaveShift = ","
	} else {
		octaveShift = "'"
	}
	for i := 0; i < octave; i++ {
		noteName += octaveShift
	}
	return noteName, nil
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ParseTabEntry parses a tab entry and returns the fret number, string number and duration
// The input format is "fret:duration\stringNum"
// duration and stringNum are optional, it uses the previous tab entry to retrieve the default duration and string
func ParseTabEntry(tabEntry string, previousStringNum int, previousDuration int) (fret int, stringNum int, duration int, err error) {
	backslashParts := strings.Split(tabEntry, "\\")
	if len(backslashParts) > 2 {
		return -1, -1, -1, fmt.Errorf("invalid input format: %s", tabEntry)
	}
	if len(backslashParts) == 2 {
		stringNum, err = strconv.Atoi(backslashParts[1])
		if err != nil {
			return -1, -1, -1, fmt.Errorf("invalid string number: %s", backslashParts[1])
		}
	} else {
		stringNum = previousStringNum
	}

	columnParts := strings.Split(backslashParts[0], ":")
	if len(backslashParts) > 2 {
		return -1, -1, -1, fmt.Errorf("invalid input format: %s", tabEntry)
	}
	if len(columnParts) == 2 {
		duration, err = strconv.Atoi(columnParts[1])
		if err != nil {
			return -1, -1, -1, fmt.Errorf("invalid string number: %s", backslashParts[1])
		}
	} else {
		duration = previousDuration
	}

	fret, err = strconv.Atoi(columnParts[0])
	if err != nil {
		return -1, -1, -1, fmt.Errorf("invalid fret number: %s", backslashParts[0])
	}
	return fret, stringNum, duration, nil
}
