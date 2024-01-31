package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	zipFileName := "denem21.zip"

	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, file := range r.File {
		fmt.Printf("Dosya Adı: %s\n", file.Name)

		// Dosyayı açın.
		fileReader, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer fileReader.Close()

		io.Copy(os.Stdout, fileReader)
		fmt.Println("\n---------------------------")
	}
}
