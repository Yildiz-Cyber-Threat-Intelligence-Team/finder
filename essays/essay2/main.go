package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// flags
	zipFilePath := flag.String("file", "", "Specify the compressed file path")
	searchText := flag.String("text", "", "Specify texts to search for by separating them with (,)")
	outputFile := flag.String("output", "", "Specify the path to the output file to save the results")
	caseSensitive := flag.Bool("case-sensitive", false, "Enable case-sensitive search")
	flag.Parse()

	// control
	if *zipFilePath == "" || *searchText == "" || *outputFile == "" {
		flag.PrintDefaults()
		return
	}

	// search words
	searchKeywords := strings.Split(*searchText, ",")

	// open
	zipFile, err := zip.OpenReader(*zipFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	// results
	var searchResults []string

	// case sensitivity
	comparisonFunc := strings.Contains
	if !*caseSensitive {
		comparisonFunc = func(s, substr string) bool {
			return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
		}
	}

	for _, file := range zipFile.File {
		if file.FileInfo().IsDir() {
			continue
		}

		// read
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(f)
		lineNumber := 1

		// scan
		for scanner.Scan() {
			line := scanner.Text()

			for _, keyword := range searchKeywords {
				if comparisonFunc(line, keyword) {
					result := fmt.Sprintf("File: %s\n", file.Name)
					result += fmt.Sprintf("Line Number: %d\n", lineNumber)
					result += fmt.Sprintf("Line: %s\n", line)
					searchResults = append(searchResults, result)
					break
				}
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		f.Close()
	}

	// output
	if len(searchResults) > 0 {
		output, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()

		for _, result := range searchResults {
			_, err := output.WriteString(result + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Search results were saved to", *outputFile)
	} else {
		fmt.Println("The searched text was not found.")
	}
}
