package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nwaples/rardecode"
)

type SearchResult struct {
	FileName string
	Lines    []string
}

func OpenArchive(path string, fileType string) (io.ReadCloser, error) {
	switch fileType {
	case ".zip":
		return os.Open(path)
	case ".rar":
		return rardecode.OpenReader(path, "")
	case ".tar.gz", ".tgz":
		fr, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return gzip.NewReader(fr)
	default:
		return nil, fmt.Errorf("Unsupported file format: %s", fileType)
	}
}

func searchInArchive(archivePath string, searchText string, fileType string, outputFilePath string) ([]SearchResult, error) {
	var results []SearchResult
	var outFile *os.File
	var err error

	startTime := time.Now()

	if outputFilePath != "" {
		outFile, err = os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		defer outFile.Close()

		if _, err := outFile.WriteString(fmt.Sprintf("Search Start Time: %s\n", startTime.Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	archive, err := OpenArchive(archivePath, fileType)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	switch archive.(type) {
	case *os.File: // ZIP file
		zipReader, err := zip.OpenReader(archivePath)
		if err != nil {
			return nil, err
		}
		defer zipReader.Close()
		for _, file := range zipReader.File {
			if !file.FileInfo().IsDir() {
				f, err := file.Open()
				if err != nil {
					return nil, err
				}
				defer f.Close()

				content, err := ioutil.ReadAll(f)
				if err != nil {
					return nil, err
				}

				lines := findLinesContainingText(string(content), searchText)
				if len(lines) > 0 {
					printResults(file.Name, searchText, lines)

					if outFile != nil {
						if err := writeResultToFile(outFile, file.Name, searchText, lines); err != nil {
							return nil, err
						}
					}

					results = append(results, SearchResult{
						FileName: file.Name,
						Lines:    lines,
					})
				}
			}
		}
	case *rardecode.ReadCloser:
		r, _ := archive.(*rardecode.ReadCloser)
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if !header.IsDir {
				content, err := ioutil.ReadAll(r)
				if err != nil {
					return nil, err
				}

				lines := findLinesContainingText(string(content), searchText)
				if len(lines) > 0 {
					printResults(header.Name, searchText, lines)

					if outFile != nil {
						if err := writeResultToFile(outFile, header.Name, searchText, lines); err != nil {
							return nil, err
						}
					}

					results = append(results, SearchResult{
						FileName: header.Name,
						Lines:    lines,
					})
				}
			}
		}
	case *gzip.Reader:
		tr := tar.NewReader(archive)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if !header.FileInfo().IsDir() {
				content, err := ioutil.ReadAll(tr)
				if err != nil {
					return nil, err
				}

				lines := findLinesContainingText(string(content), searchText)
				if len(lines) > 0 {
					printResults(header.Name, searchText, lines)

					if outFile != nil {
						if err := writeResultToFile(outFile, header.Name, searchText, lines); err != nil {
							return nil, err
						}
					}

					results = append(results, SearchResult{
						FileName: header.Name,
						Lines:    lines,
					})
				}
			}
		}
	default:
		return nil, fmt.Errorf("Unsupported archive type")
	}

	duration := time.Since(startTime)

	if outFile != nil {
		if _, err := outFile.WriteString(fmt.Sprintf("Search End Time: %s\n", time.Now().Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	fmt.Println("Total time:", duration)

	return results, nil
}

func searchInTextFile(filePath string, searchText string, outputFilePath string, deleteFile bool) ([]SearchResult, error) {
	var results []SearchResult
	var outFile *os.File
	var err error

	startTime := time.Now()

	if outputFilePath != "" {
		outFile, err = os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		defer outFile.Close()

		if _, err := outFile.WriteString(fmt.Sprintf("Search Start Time: %s\n", startTime.Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := findLinesContainingText(string(fileContent), searchText)
	if len(lines) > 0 {
		printResults(filePath, searchText, lines)

		if outFile != nil {
			if err := writeResultToFile(outFile, filePath, searchText, lines); err != nil {
				return nil, err
			}
		}

		results = append(results, SearchResult{
			FileName: filePath,
			Lines:    lines,
		})
	}

	duration := time.Since(startTime)

	if outFile != nil {
		if _, err := outFile.WriteString(fmt.Sprintf("Search End Time: %s\n", time.Now().Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	fmt.Println("Total time:", duration)

	if deleteFile {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File successfully deleted.")
		}
	}

	return results, nil
}

func findLinesContainingText(content string, searchText string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchText) {
			lines = append(lines, line)
		}
	}
	return lines
}

func printResults(fileName string, searchText string, lines []string) {
	fmt.Println("File:", fileName)
	fmt.Println("Searched Text:", searchText)
	fmt.Println("Lines:")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("---------------------------------------")
}

func writeResultToFile(file *os.File, fileName string, searchText string, lines []string) error {

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString(fmt.Sprintf("File: %s\n", fileName)); err != nil {
		return err
	}

	if _, err := writer.WriteString(fmt.Sprintf("Search Text: %s\n", searchText)); err != nil {
		return err
	}

	if _, err := writer.WriteString("Lines:\n"); err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
	}

	if _, err := writer.WriteString("---------------------------------------\n"); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println(`
Usage Guide:

    Search for text within an archive file:
        main <archive_file> "<search_text>" <file_type> [output_file] [-d]
            <archive_file>: The path to the archive file to search within.
            "<search_text>": The text to search for (enclose in double quotes if it contains spaces).
            <file_type>: The type of the archive file (e.g., .zip, .rar, .tar.gz).
            [output_file]: Optional. The path to the output file where search results will be saved.
            [-d]: Optional. Flag to delete the original file after searching.

    Search for text within a text file:
        main <text_file> "<search_text>" [output_file] [-d]
            <text_file>: The path to the text file to search within.
            "<search_text>": The text to search for (enclose in double quotes if it contains spaces).
            [output_file]: Optional. The path to the output file where search results will be saved.
            [-d]: Optional. Flag to delete the original file after searching.
    `)
}

func main() {
	args := os.Args

	if len(args) < 4 || len(args) > 6 {
		printUsage()
		return
	}

	var archiveFilePath, searchText, fileType, outputFilePath string
	var deleteFile bool

	if len(args) >= 5 {
		for _, arg := range args[4:] {
			if arg == "-d" {
				deleteFile = true
			} else {
				outputFilePath = arg
			}
		}
	}

	archiveFilePath = args[1]
	searchText = args[2]
	fileType = args[3]

	fullArchivePath := filepath.Clean(archiveFilePath)

	searchText = strings.Trim(searchText, "\"")

	if fileType != ".zip" && fileType != ".rar" && fileType != ".tar.gz" && fileType != ".tgz" {
		results, err := searchInTextFile(fullArchivePath, searchText, outputFilePath, deleteFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("Search completed. No results found.")
			return
		}

		fmt.Println("Search completed. Results printed to the terminal.")

		if outputFilePath != "" {
			fmt.Println("Results saved to", outputFilePath)
		}

		return
	}

	results, err := searchInArchive(fullArchivePath, searchText, fileType, outputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("Search completed. No results found.")
		return
	}

	fmt.Println("Search completed. Results printed to the terminal.")

	if outputFilePath != "" {
		fmt.Println("Results saved to", outputFilePath)
	}

	if deleteFile {
		err := os.Remove(fullArchivePath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File successfully deleted.")
		}
	}
}
