package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func main() {
	input := flag.String("input", "", "File name to read contents and send")
	url := flag.String("url", "", "URL to send contents to")
	fieldName := flag.String("field-name", "files", "Field name to append files to")

	flag.Parse()

	filePaths := strings.Split(*input, ",")

	if len(filePaths) == 0 {
		fmt.Println("Input file name must be specified")

		return
	}

	if len(*url) == 0 {
		fmt.Println("URL file name must be specified")

		return
	}

	if len(*fieldName) == 0 {
		fmt.Println("Field name file name must be specified")

		return
	}

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for _, filePath := range filePaths {
		err := addFileToWriter(writer, filePath, *fieldName)

		if err != nil {
			fmt.Println("Error adding file:", err)

			return
		}
	}

	err := writer.Close()

	if err != nil {
		fmt.Println("Error closing writer:", err)

		return
	}

	resp, err := http.Post(*url, writer.FormDataContentType(), &body)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func addFileToWriter(writer *multipart.Writer, filePath, fieldName string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filePath)

	if err != nil {
		return fmt.Errorf("failed to create form field for file %s: %w", filePath, err)
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return fmt.Errorf("failed to copy file content for file %s: %w", filePath, err)
	}

	return nil
}
