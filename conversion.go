package main

import (
	"fmt"
	"log/slog"
	"strconv"
)

// Exported standard tuning for guitar (E A D G B e)
var GuitarTuning = []string{"E", "A", "D", "G", "B", "e"}

// Exported standard tuning for bass (E A D G)
var BassTuning = []string{"E", "A", "D", "G"}

func convertToLilypond(fret, stringNum int, tuning []string, duration string) string {
	// Note names in LilyPond format
	noteNames := []string{"c", "cis", "d", "dis", "e", "f", "fis", "g", "gis", "a", "ais", "b"}

	slog.Debug("Convert tab entry to LilyPond format", "fret", fret, "stringNum", stringNum, "tuning", tuning, "duration", duration)

	// Calculate the note name and octave
	openStringNote := tuning[stringNum-1]
	openStringNoteIndex := getNoteIndex(openStringNote)
	noteIndex := (openStringNoteIndex + fret) % 12
	octave := (openStringNoteIndex + fret) / 12

	slog.Debug("Calculate the note name and octave", "openStringNote", openStringNote, "openStringNoteIndex", openStringNoteIndex, "noteIndex", noteIndex, "octave", octave)

	// Adjust octave for LilyPond notation
	if noteIndex < openStringNoteIndex {
		octave++
	}

	// Convert duration to an integer to remove leading zeros
	durationInt, _ := strconv.Atoi(duration)
	return fmt.Sprintf("%s'%d%d\\%d", noteNames[noteIndex], octave, durationInt, stringNum)
}

func getNoteIndex(note string) int {
	switch note {
	case "C":
		return 0
	case "C#":
		return 1
	case "D":
		return 2
	case "D#":
		return 3
	case "E":
		return 4
	case "F":
		return 5
	case "F#":
		return 6
	case "G":
		return 7
	case "G#":
		return 8
	case "A":
		return 9
	case "A#":
		return 10
	case "B":
		return 11
	default:
		return -1
	}
}
