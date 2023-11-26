package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// countWords reads the content of a file and utilizes string.Fields
// to count the number of words.
func countWords(fileName string) (int, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	contentString := string(content)
	return len(strings.Fields(contentString)), nil
}

// countLinesAndBytes reads a file line by line using a Scanner
// and returns the number of lines and total bytes.
func countLinesAndBytes(fileName string) (int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var lines, bytes int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
		bytes += len(scanner.Bytes()) + 1 // +1 for the newline character
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return bytes, lines, nil
}

func main() {
	app := &cli.App{
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			// Specify command-line options for counting
			// bytes, lines, and words
			&cli.BoolFlag{
				Name:    "bytes",
				Aliases: []string{"c"},
				Usage:   "print the byte counts",
			},
			&cli.BoolFlag{
				Name:    "lines",
				Aliases: []string{"l"},
				Usage:   "print the newline counts",
			},
			&cli.BoolFlag{
				Name:    "words",
				Aliases: []string{"w"},
				Usage:   "print the word counts",
			},
		},
		Action: func(cliContext *cli.Context) error {
			var counts []int
			var fileName string
			var numBytes, numWords, numLines int
			var err error

			// Check if a filename is provided
			if cliContext.NArg() == 0 {
				return fmt.Errorf("missing a filename")
			}
			fileName = cliContext.Args().Get(0)

			// Check if no options are provided, which is the equivalent to the -clw
			if cliContext.NumFlags() == 0 {
				numBytes, numLines, err = countLinesAndBytes(fileName)
				if err != nil {
					return err
				}
				numWords, err = countWords(fileName)
				if err != nil {
					return err
				}
				counts = append(counts, numLines, numWords, numBytes)
			}

			// Count words if the corresponding flag is set
			if cliContext.Bool("words") {
				numWords, err = countWords(fileName)
				if err != nil {
					return err
				}
			}

			// Count lines and bytes if the corresponding flags are set
			if cliContext.Bool("lines") || cliContext.Bool("bytes") {
				numBytes, numLines, err = countLinesAndBytes(fileName)
				if err != nil {
					return err
				}
			}

			// Collect the counts based on the flags
			if cliContext.Bool("lines") {
				counts = append(counts, numLines)
			}
			if cliContext.Bool("words") {
				counts = append(counts, numWords)
			}
			if cliContext.Bool("bytes") {
				counts = append(counts, numBytes)
			}

			// Convert counts to strings and concatenate them
			countsStr := make([]string, len(counts))
			for i, num := range counts {
				countsStr[i] = fmt.Sprint(num)
			}
			result := strings.Join(countsStr, " ")

			// Print the result
			fmt.Printf("%v %v\n", result, fileName)
			return nil
		},
	}

	// Run the CLI application, and check for any errors during execution.
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
