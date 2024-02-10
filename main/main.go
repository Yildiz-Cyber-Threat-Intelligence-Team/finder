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

	"github.com/fatih/color"
	"github.com/nwaples/rardecode"
)

type SearchResult struct {
	FileName string
	Lines    []string // Satırlardaki metinlerin listesi
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
		return nil, fmt.Errorf("Dosya formatı desteklenmiyor: %s", fileType)
	}
}

func searchInArchive(archivePath string, searchText string, fileType string, outputFilePath string) ([]SearchResult, error) {
	var results []SearchResult
	var outFile *os.File
	var err error

	// Başlangıç zamanını al
	startTime := time.Now()

	// Eğer bir çıktı dosyası belirtilmişse, dosyayı üstüne yazacak şekilde aç
	if outputFilePath != "" {
		outFile, err = os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		defer outFile.Close()

		// Dosyaya başlama zamanını yaz
		if _, err := outFile.WriteString(fmt.Sprintf("Arama Başlangıç Zamanı: %s\n", startTime.Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	archive, err := OpenArchive(archivePath, fileType)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	switch archive.(type) {
	case *os.File:
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
		return nil, fmt.Errorf("Desteklenmeyen arşiv türü")
	}

	// Geçen zamanı hesapla
	duration := time.Since(startTime)

	// Sonuç dosyasına bitiş zamanını yaz
	if outFile != nil {
		if _, err := outFile.WriteString(fmt.Sprintf("Arama Bitiş Zamanı: %s\n", time.Now().Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	fmt.Println("Toplam süre:", duration)

	return results, nil
}

// searchInTextFile, normal bir metin dosyası içinde arama yapar.
func searchInTextFile(filePath string, searchText string, outputFilePath string) ([]SearchResult, error) {
	var results []SearchResult
	var outFile *os.File
	var err error

	// Başlangıç zamanını al
	startTime := time.Now()

	// Eğer bir çıktı dosyası belirtilmişse, dosyayı üstüne yazacak şekilde aç
	if outputFilePath != "" {
		outFile, err = os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		defer outFile.Close()

		// Dosyaya başlama zamanını yaz
		if _, err := outFile.WriteString(fmt.Sprintf("Arama Başlangıç Zamanı: %s\n", startTime.Format("02.01.2006 15:04"))); err != nil {
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

	// Geçen zamanı hesapla
	duration := time.Since(startTime)

	// Sonuç dosyasına bitiş zamanını yaz
	if outFile != nil {
		if _, err := outFile.WriteString(fmt.Sprintf("Arama Bitiş Zamanı: %s\n", time.Now().Format("02.01.2006 15:04"))); err != nil {
			return nil, err
		}
	}

	fmt.Println("Toplam süre:", duration)

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
	fmt.Println("Dosya:", fileName)
	fmt.Println("Aranan Metin:", searchText)
	fmt.Println("Satırlar:")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("---------------------------------------")
}

func writeResultToFile(file *os.File, fileName string, searchText string, lines []string) error {
	// Dosya bilgilerini dosyaya yaz
	if _, err := file.WriteString(fmt.Sprintf("Dosya: %s\n", fileName)); err != nil {
		return err
	}
	// Aranan metni dosyaya yaz
	if _, err := file.WriteString(fmt.Sprintf("Aranan Metin: %s\n", searchText)); err != nil {
		return err
	}
	// Eşleşen satırları dosyaya yaz
	if _, err := file.WriteString("Satırlar:\n"); err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
	}
	// Ayırıcı satırı dosyaya yaz
	if _, err := file.WriteString("---------------------------------------\n"); err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println(`
Arama Aracı Kullanım Kılavuzu:

Kullanım: 
  go run main.go "arsiv_dosya_yolu" "arama_metni" "dosya_turu" [cikti_dosyasi]

Örnek:
  go run main.go "dosya.rar" "searchWord" ".rar" [output.txt]

Parametreler:
  - "arsiv_dosya_yolu": Arama yapılacak arşiv dosyasının veya metin dosyasının yolu.
  - "arama_metni": Arşiv dosyası veya metin dosyası içinde aranacak metin.
  - "dosya_turu": Arşiv dosyasının türü (".zip", ".rar", ".tar.gz", ".tgz"). Metin dosyası için bu parametre kullanılmaz.
  - [cikti_dosyasi]: İsteğe bağlı olarak, sonuçların kaydedileceği dosyanın yolu.

Örnekler:
  - Dosya.rar arşiv dosyasında "searchWord" ifadesini aramak için:
    go run main.go "dosya.rar" "searchWord" ".rar"
  - Metin.txt dosyasında "arama kelimesi" ifadesini arayıp sonuçları output.txt dosyasına kaydetmek için:
    go run main.go "metin.txt" "arama kelimesi" output.txt
	`)
}

func displayAsciiArt() {
	file, err := os.ReadFile("ascii_art.txt")
	if err != nil {
		c := color.New(color.FgRed, color.Bold)
		c.Println("ASCII art could not be displayed:", err)
		return
	}
	fmt.Println(string(file))
}

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		printUsage()
		return
	}

	archiveFilePath := os.Args[1]
	searchText := os.Args[2]
	fileType := os.Args[3]
	var outputFilePath string
	if len(os.Args) == 5 {
		outputFilePath = os.Args[4]
	}

	// Dosya yolunu tam olarak almak için filepath.Clean kullanın
	fullArchivePath := filepath.Clean(archiveFilePath)

	// Arama metnini tırnak içine alın
	searchText = strings.Trim(searchText, "\"")

	// Dosya türü arşiv dosyası değilse, metin dosyası içinde arama yap
	if fileType != ".zip" && fileType != ".rar" && fileType != ".tar.gz" && fileType != ".tgz" {
		results, err := searchInTextFile(fullArchivePath, searchText, outputFilePath)
		if err != nil {
			fmt.Println("Hata:", err)
			return
		}

		// Eğer sonuç yoksa dosyaya yazma işlemi yapma
		if len(results) == 0 {
			fmt.Println("Arama tamamlandı. Sonuç bulunamadı.")
			return
		}

		fmt.Println("Arama tamamlandı. Sonuçlar terminale yazıldı.")

		// Dosyaya yazma işlemi yapılmışsa bilgi mesajı göster
		if outputFilePath != "" {
			fmt.Println("Sonuçlar", outputFilePath, "dosyasına kaydedildi.")
		}

		return
	}

	// Arşiv dosyası içinde arama yap
	results, err := searchInArchive(fullArchivePath, searchText, fileType, outputFilePath)
	if err != nil {
		fmt.Println("Hata:", err)
		return
	}

	// Eğer sonuç yoksa dosyaya yazma işlemi yapma
	if len(results) == 0 {
		fmt.Println("Arama tamamlandı. Sonuç bulunamadı.")
		return
	}

	fmt.Println("Arama tamamlandı. Sonuçlar terminale yazıldı.")

	// Dosyaya yazma işlemi yapılmışsa bilgi mesajı göster
	if outputFilePath != "" {
		fmt.Println("Sonuçlar", outputFilePath, "dosyasına kaydedildi.")
	}
}
