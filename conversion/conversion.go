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

var Newline = "-newline-"

var second = 2
var sixteenth = 16

type TmplData struct {
	TablilyNotes string
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
func ConvertToLilypond(fret, stringNum int, tuning []string, duration int, deadNote bool, rest bool, currentOctave int) (string, int, error) {

	slog.Debug("Convert tab entry to LilyPond format", "fret", fret, "stringNum", stringNum, "tuning", tuning, "duration", duration)

	if rest {
		return fmt.Sprintf("r%d", duration), currentOctave, nil
	}

	// Calculate the note name and octave
	openStringNote, err := getOpenStringNote(tuning, stringNum)
	if err != nil {
		slog.Error("Invalid string number", "stringNum", stringNum)
		return "", 0, err
	}
	openStringNoteIndex, err := getNoteIndex(openStringNote)
	if err != nil {
		slog.Error("Invalid open string note", "openStringNote", openStringNote)
		return "", 0, err
	}
	noteIndex := openStringNoteIndex + fret

	slog.Debug("Calculate the note index", "openStringNote", openStringNote, "openStringNoteIndex", openStringNoteIndex, "noteIndex", noteIndex)

	noteName, currentOctave, err := findNoteName(noteIndex, currentOctave)
	if err != nil {
		slog.Error("Invalid note index", "noteIndex", noteIndex)
		return "", 0, err
	}

	// Convert duration to an integer to remove leading zeros
	return fmt.Sprintf("%s%d\\%d", noteName, duration, stringNum), currentOctave, nil
}

func findNoteName(noteIndex int, currentOctave int) (string, int, error) {
	var octave int
	var octaveShift int
	noteName := ""

	scaleIndex := mod(noteIndex, 12)
	slog.Debug("Find note name", "noteIndex", noteIndex, "scaleIndex", scaleIndex)
	for note, i := range noteMap {
		if i == scaleIndex {
			noteName = note
		}
	}
	if noteName == "" {
		return "", 0, fmt.Errorf("invalid note index: %d", noteIndex)
	}

	if noteIndex >= 0 {
		octave = noteIndex / 12
	} else {
		octave = Abs(noteIndex)/12 - 1
	}

	slog.Debug("Calculate the octave", "octave", octave, "currentOctave", currentOctave)

	octaveShift = octave - currentOctave

	if octaveShift < 0 {
		for i := 0; i < Abs(octaveShift); i++ {
			noteName += ","
		}
	} else if octaveShift > 0 {
		for i := 0; i < Abs(octaveShift); i++ {
			noteName += "'"
		}
	}
	return noteName, octave, nil
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
// TODO manage 'r' for rest
func ParseTabEntry(tabEntry string, previousStringNum int, previousDuration int, tuning []string) (fret int, stringNum int, duration int, deadnote bool, rest bool, err error) {
	backslashParts := strings.Split(tabEntry, "\\")
	slog.Debug("Parse tab entry", "tabEntry", tabEntry, "backslashParts", backslashParts, "previousStringNum", previousStringNum, "previousDuration", previousDuration)
	if len(backslashParts) > 2 {
		return -1, -1, -1, false, false, fmt.Errorf("invalid input format: %s", tabEntry)
	}
	if len(backslashParts) == 2 {
		stringNum, err = strconv.Atoi(backslashParts[1])
		if err != nil || stringNum < 1 || stringNum > len(tuning) {
			return -1, -1, -1, false, false, fmt.Errorf("invalid string number: %s in note %s", backslashParts[1], tabEntry)
		}
	} else {
		stringNum = previousStringNum
	}

	colonParts := strings.Split(backslashParts[0], ":")
	if len(backslashParts) > 2 {
		return -1, -1, -1, false, false, fmt.Errorf("invalid input format: %s", tabEntry)
	}
	if len(colonParts) == 2 {
		duration, err = strconv.Atoi(colonParts[1])
		if err != nil {
			return -1, -1, -1, false, false, fmt.Errorf("invalid duration: %s in note %s", colonParts[1], tabEntry)
		}
	} else {
		duration = previousDuration
	}

	if colonParts[0] == "x" {
		fret = 0
		deadnote = true
	} else if colonParts[0] == "r" {
		fret = 0
		rest = true
	} else {
		fret, err = strconv.Atoi(colonParts[0])
		if err != nil {
			return -1, -1, -1, false, false, fmt.Errorf("invalid fret number: %s in note %s", backslashParts[0], tabEntry)
		}
	}
	return fret, stringNum, duration, deadnote, rest, nil
}

func mod(d int, m int) int {
	d %= m
	if d < 0 {
		d += m
	}
	return d
}
