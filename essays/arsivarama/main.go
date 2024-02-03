package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func searchInZip(zipFilePath string, searchText string, outputFilePath string) error {
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	for _, file := range zipFile.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		lineNumber := 0

		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if strings.Contains(line, searchText) {
				result := fmt.Sprintf("File: %s, Line: %d, Content: %s", file.Name, lineNumber, line)
				fmt.Println(result)
				writer.WriteString(result + "\n")
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func main() {
	fmt.Println("1. Arşiv Dosyasının Yolunu Girin:")
	var zipFilePath string
	fmt.Scanln(&zipFilePath)

	fmt.Println("2. Aranacak Metni veya Aranacak Kelimelerin Bulunduğu Dosyanın Yolunu Girin:")
	var searchText string
	fmt.Scanln(&searchText)

	fmt.Println("3. Çıktı Dosyasının Yolunu Girin:")
	var outputFilePath string
	fmt.Scanln(&outputFilePath)

	err := searchInZip(zipFilePath, searchText, outputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
