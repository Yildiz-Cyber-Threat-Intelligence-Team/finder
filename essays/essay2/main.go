package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	displayAsciiArt()

	// flags
	zipFilePath := flag.String("file", "", "Specify the compressed file path - required")
	searchText := flag.String("text", "", "Specify the text to search for or specify multiple texts separated by (,) - required")
	outputFile := flag.String("output", "", "Specify the path to the output file to save the results - optional")
	caseSensitive := flag.Bool("case-sensitive", false, "Specify whether the search should be case-sensitive - optional")
	helpFlag := flag.Bool("help", false, "Help for using the finder tool")
	deleteFile := flag.Bool("del", false, "Specify whether to delete the compressed file after search - optional")
	flag.Parse()

	// control
	if *helpFlag {
		fmt.Println("Flags")
		flag.PrintDefaults()
		return
	}

	if *zipFilePath == "" || *searchText == "" {
		color.Red("Please run the -help command to use the Finder tool.")
		return
	}

	color.Green("The Finder Tool is starting...\n")

	// separate words
	searchKeywords := strings.Split(*searchText, ",")

	// open
	zipFile, err := zip.OpenReader(*zipFilePath)
	if err != nil {
		color.Red("The compressed file was not found.")
		return
	}
	defer zipFile.Close()

	// results
	var searchResults []string

	var totalSize int64
	fileCount := 0

	for _, file := range zipFile.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fileCount++

		// increment total size
		totalSize += file.FileInfo().Size()

		// read
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close() // close file after search

		scanner := bufio.NewScanner(f)
		lineNumber := 1

		// scan
		for scanner.Scan() {
			line := scanner.Text()
			compareLine := line
			if !*caseSensitive {
				compareLine = strings.ToLower(line)
			}

			for _, keyword := range searchKeywords {
				compareKeyword := keyword
				if !*caseSensitive {
					compareKeyword = strings.ToLower(keyword)
				}

				if strings.Contains(compareLine, compareKeyword) {
					result := fmt.Sprintf("File: %s\n", file.Name)
					result += fmt.Sprintf("Line Number: %d\n", lineNumber)
					result += fmt.Sprintf("Line: %s\n", line)
					result += "--------------------\n"
					searchResults = append(searchResults, result)
					break
				}
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// total size and file count
	fmt.Printf("Compressed File Properties\n")
	fmt.Printf("Total Size: %d bytes\n", totalSize)
	fmt.Printf("File Count: %d\n", fileCount)
	fmt.Printf("--------------------\n\n")

	// output
	if len(searchResults) > 0 {
		if *outputFile != "" {
			output, err := os.Create(*outputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			for _, result := range searchResults {
				_, err := output.WriteString(result)
				if err != nil {
					log.Fatal(err)
				}
			}

			color.Green("Search results were saved to %s", *outputFile)
		} else {
			for _, result := range searchResults {
				fmt.Println(result)
			}
			color.Green("Search results printed to the terminal.")
		}
	} else {
		color.Red("The searched text was not found.")
	}

	if *deleteFile {
		zipFile.Close() // close file
		err := os.Remove(*zipFilePath)
		if err != nil {
			color.Red("The compressed file has not been deleted.")
		}
		color.Blue("The compressed file has been deleted.")
	}

}
func displayAsciiArt() {
	file, err := os.ReadFile("ascii_art.txt")
	if err != nil {
		color.Red("ASCII art could not be displayed.")
		return
	}
	fmt.Println(string(file))
}
