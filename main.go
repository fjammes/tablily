package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fjammes/tablily/conversion"
	"github.com/fjammes/tablily/log"
)

func main() {
	// Define command-line flags
	instrument := flag.String("instrument", "guitar", "Instrument type (guitar/bass)")
	inputFile := flag.String("input", "", "Input file containing tab entries")
	verbosity := flag.Int("v", 1, "Log level (0-4)")
	flag.Parse()

	log.Init(*verbosity)

	// Validate instrument type
	var tuning []string
	if *instrument == "guitar" {
		tuning = conversion.GuitarTuning
	} else if *instrument == "bass" {
		tuning = conversion.BassTuning
	} else {
		fmt.Println("Invalid instrument type")
		return
	}

	var input string
	if *inputFile != "" {
		// Read input from file
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			input += scanner.Text() + " "
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		input = strings.TrimSpace(input)
	} else {
		// Read input from stdin
		fmt.Println("Enter tab (format: fret[:duration][/string], e.g., 3/1/4 for 3rd fret on 1st string with duration 4, r for rest, x/string for ghost note):")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
	}

	slog.Debug("Retrieve the input", "input", input)
	notes := strings.Split(input, " ")
	lilypondNotes := []string{}

	previousStringNum := -1
	previousDuration := 4
	for _, note := range notes {
		// Parse the tab entry
		fret, stringNum, duration, deadNote, rest, err := conversion.ParseTabEntry(note, previousStringNum, previousDuration, tuning)
		if err != nil {
			fmt.Println("Error parsing tab entry:", err)
			os.Exit(1)
		}

		// TODO manage previous octave, string and duration
		lilypondNote, err := conversion.ConvertToLilypond(fret, stringNum, tuning, duration, deadNote, rest)
		if err != nil {
			fmt.Println("Error converting tab entry to LilyPond format:", err)
			os.Exit(1)
		}
		lilypondNotes = append(lilypondNotes, lilypondNote)
	}

	// Output the result in LilyPond format
	fmt.Println("LilyPond format:")
	fmt.Println(strings.Join(lilypondNotes, " "))
}
