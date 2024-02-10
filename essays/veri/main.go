package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Sıkıştırılmış dosyanın adını belirtin.
	zipFileName := "denem21.zip"

	// Hedef dizini belirtin. Bu dizinde belirli dosyaları kaydedeceğiz.
	destDirectory := "veri"

	// Hangi dosyaları kopyalamak istediğinizi belirtin (örneğin, ".txt" uzantılı dosyaları kopyalamak için).
	targetFileExtensions := []string{".txt"}

	// Zip dosyasını açın.
	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		fmt.Println("Sıkıştırılmış dosya açılamadı:", err)
		return
	}
	defer r.Close()

	// Hedef dizini oluşturun (eğer yoksa).
	err = os.MkdirAll(destDirectory, os.ModePerm)
	if err != nil {
		fmt.Println("Hedef dizin oluşturulamadı:", err)
		return
	}

	// Zip dosyasındaki her dosyayı dolaşın.
	for _, file := range r.File {
		// Dosya adını alın.
		fileName := filepath.Join(destDirectory, filepath.Base(file.Name))

		// Hedef dizine kopyalamak istediğiniz dosyaları belirleyin.
		if shouldCopyFile(fileName, targetFileExtensions) {
			// Dosyayı açın.
			fileReader, err := file.Open()
			if err != nil {
				fmt.Println("Dosya açılamadı:", err)
				continue
			}
			defer fileReader.Close()

			// Dosyayı hedef dizine kaydedin.
			destFile, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Dosya kaydedilemedi:", err)
				continue
			}

			// Dosyayı kopyalayın.
			_, err = io.Copy(destFile, fileReader)
			destFile.Close()
			if err != nil {
				fmt.Println("Dosya kopyalanamadı:", err)
				continue
			}

			fmt.Printf("Dosya kaydedildi: %s\n", fileName)
		}
	}

	// Zip dosyasını silin.
	err = os.Remove(zipFileName)
	if err != nil {
		fmt.Println("Sıkıştırılmış dosya silinemedi:", err)
		return
	}

	fmt.Println("Sıkıştırılmış dosya başarıyla işlendi ve silindi.")
}

// Belirli dosyanın kopyalanıp kopyalanmayacağını kontrol eden yardımcı fonksiyon.
func shouldCopyFile(fileName string, targetFileExtensions []string) bool {
	ext := filepath.Ext(fileName)
	for _, targetExt := range targetFileExtensions {
		if ext == targetExt {
			return true
		}
	}
	return false
}
