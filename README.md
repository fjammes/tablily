# Tablature to LilyPond Converter

This Go program converts guitar and bass tablature input into LilyPond format. It supports standard tuning for both guitar and bass, and can handle rests, ghost notes, and note durations.

## Features

- Converts guitar and bass tablature to LilyPond format.
- Supports standard tuning for guitar (E A D G B e) and bass (E A D G).
- Handles rests (r), ghost notes (x/string), and note durations (fret/string/duration).

## Usage

Command-Line Options
-instrument: Specify the instrument type (guitar or bass). Default is guitar.
-input: Specify the input file containing tab entries. If not provided, the program will read from standard input.

### Input Format

- fret/string[/duration]: Specifies the fret and string number, with an optional duration. For example, 3/1/4 means the 3rd fret on the 1st string with a duration of 4.
- r: Represents a rest (silence) with a default duration of 4.
- x/string: Represents a ghost note on the specified string with a default duration of 4.

Example tab.txt File
```
3/1/4 2/2/8 r 0/3/4 x/4/4 1/4/2 3/5/4 2/6/4
5/1/4 4/2/8 r 2/3/4 x/3/4 3/4/2 5/5/4 4/6/4
7/1/4 6/2/8 r 4/3/4 x/2/4 5/4/2 7/5/4 6/6/4
```


## Installation
Clone the repository:

```sh
git clone https://github.com/fjammes/tablily.git
cd tablily
go build
# run unit tests
go test
```

## Run the Program
To run the program with an input file:

```sh
tablily -instrument=guitar -input=tab.txt
```

To run the program and enter tab input manually:

```sh
tablily -instrument=bass
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgements
LilyPond for the music notation format.
The Go programming language.