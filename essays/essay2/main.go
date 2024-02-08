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
	zipFilePath := flag.String("file", "", "Specify the compressed file path - required")
	searchText := flag.String("text", "", "Specify the text to search for or specify multiple texts separated by (,) - required")
	outputFile := flag.String("output", "", "Specify the path to the output file to save the results - optional")
	caseSensitive := flag.Bool("case-sensitive", false, "Specify whether the search should be case-sensitive - optional")
	helpFlag := flag.Bool("help", false, "Help for using the finder tool")
	flag.Parse()

	// control
	if *helpFlag {
		fmt.Println("Flags;")
		flag.PrintDefaults()
		return
	}

	if *zipFilePath == "" || *searchText == "" {
		fmt.Println("Please run the -help command to use the Finder tool.")
		return
	}

	// separate words
	searchKeywords := strings.Split(*searchText, ",")

	// open
	zipFile, err := zip.OpenReader(*zipFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	// results
	var searchResults []string

	//
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
			compareLine := line
			if !*caseSensitive {
				compareLine = strings.ToLower(line)
			}

			//compare words
			for _, keyword := range searchKeywords {
				compareKeyword := keyword
				if !*caseSensitive {
					compareKeyword = strings.ToLower(keyword)
				}

				if strings.Contains(compareLine, compareKeyword) {
					result := fmt.Sprintf("File: %s\n", file.Name)
					result += fmt.Sprintf("Line Number: %d\n", lineNumber)
					result += fmt.Sprintf("Line: %s\n", line)
					result += "----------\n"
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

			fmt.Println("Search results were saved to", *outputFile)
		} else {
			for _, result := range searchResults {
				fmt.Println(result)
			}
		}
	} else {
		fmt.Println("The searched text was not found.")
	}
}
