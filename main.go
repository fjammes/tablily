package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/fjammes/tablily/log"
)

// Standard tuning for guitar (E A D G B e)
var guitarTuning = []string{"E", "A", "D", "G", "B", "e"}

// Standard tuning for bass (E A D G)
var bassTuning = []string{"E", "A", "D", "G"}

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
		tuning = guitarTuning
	} else if *instrument == "bass" {
		tuning = bassTuning
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
		fmt.Println("Enter tab (format: fret/string[/duration], e.g., 3/1/4 for 3rd fret on 1st string with duration 4, r for rest, x/string for ghost note):")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
	}

	slog.Debug("Retrieve the input", "input", input)
	notes := strings.Split(input, " ")
	lilypondNotes := []string{}

	for _, note := range notes {
		if note == "r" {
			lilypondNotes = append(lilypondNotes, "r4")
			continue
		}

		parts := strings.Split(note, "\\")
		if len(parts) < 2 || len(parts) > 3 {
			fmt.Println("Invalid input format")
			return
		}

		if parts[0] == "x" {
			stringNum, err := strconv.Atoi(parts[1])
			if err != nil || stringNum < 1 || stringNum > len(tuning) {
				fmt.Println("Invalid string number for ghost note")
				return
			}
			duration := "4"
			if len(parts) == 3 {
				duration = parts[2]
			}
			lilypondNotes = append(lilypondNotes, fmt.Sprintf("/deadNote x%s", duration))
			continue
		}

		fret, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Invalid fret number")
			return
		}
		slog.Debug("Retrieve the fret number", "fret", fret)

		stringNum, err := strconv.Atoi(parts[1])
		if err != nil || stringNum < 1 || stringNum > len(tuning) {
			fmt.Println("Invalid string number")
			return
		}

		duration := "4"
		if len(parts) == 3 {
			duration = parts[2]
		}

		lilypondNote := convertToLilypond(fret, stringNum, tuning, duration)
		lilypondNotes = append(lilypondNotes, lilypondNote)
	}

	// Output the result in LilyPond format
	fmt.Println("LilyPond format:")
	fmt.Println(strings.Join(lilypondNotes, " "))
}
