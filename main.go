package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/fjammes/tablily/conversion"
	"github.com/fjammes/tablily/log"
)

func main() {
	// Define command-line flags
	instrument := flag.String("I", "guitar", "Instrument type (guitar/bass)")
	inputFile := flag.String("i", "", "Input file containing tab entries")
	tmplFile := flag.String("t", "", "Ouput file: template lilypond file")
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
	var notes []string
	if *inputFile != "" {
		// Read input from file
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		// Read the file line by line
		// Store \n in input to keep the same format
		for scanner.Scan() {
			line := scanner.Text()

			lineTokens := strings.Fields(line)
			lineTokens = append(lineTokens, conversion.Newline)
			// Append the tokens to the result slice
			notes = append(notes, lineTokens...)
			slog.Debug("Read line", "line", line)
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
		notes = strings.Split(input, " ")
	}

	slog.Debug("Retrieve the input", "input", input)
	// split the input into notes on multiple white spaces

	lilypondNotes := []string{}

	previousStringNum := -1
	previousDuration := 4
	octave := 0
	var lilypondNote string
	for _, note := range notes {
		if note != conversion.Newline {
			// Parse the tab entry
			fret, stringNum, duration, deadNote, rest, err := conversion.ParseTabEntry(note, previousStringNum, previousDuration, tuning)
			if err != nil {
				fmt.Println("Error parsing tab entry:", err)
				os.Exit(1)
			}
			previousStringNum = stringNum
			previousDuration = duration

			// TODO manage previous octave, string and duration
			lilypondNote, octave, err = conversion.ConvertToLilypond(fret, stringNum, tuning, duration, deadNote, rest, octave)
			slog.Debug("Convert to LilyPond", "lilypondNote", lilypondNote, "octave", octave)
			if err != nil {
				fmt.Println("Error converting tab entry to LilyPond format:", err)
				os.Exit(1)
			}
			lilypondNotes = append(lilypondNotes, lilypondNote)
		} else {
			lilypondNotes = append(lilypondNotes, "\n")
		}
	}

	// Output the result in LilyPond format
	data := conversion.TmplData{
		TablilyNotes: strings.Join(lilypondNotes, " "),
	}

	// Open output file and replace {.tablily_notes} with the LilyPond notes
	// Use go template to replace the placeholder
	name := path.Base(*tmplFile)
	tmpl, err := template.New(name).ParseFiles(*tmplFile)
	if err != nil {
		slog.Error("Error parsing template file", "error", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		slog.Error("Error executing template", "error", err)
		os.Exit(1)
	}
	fmt.Println(buf.String())
}
