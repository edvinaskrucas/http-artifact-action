package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	inputFlag := flag.String("input", "", "File name to read contents and send")
	urlFlag := flag.String("url", "", "URL to send contents to")
	fieldNameFlag := flag.String("field-name", "files", "Field name to append files to")
	dataFlag := flag.String("data", "", "Extra data to send")

	flag.Parse()

	filePaths := strings.Split(*inputFlag, ",")

	if len(filePaths) == 0 {
		fmt.Println("Input file name must be specified")

		return
	}

	if len(*urlFlag) == 0 {
		fmt.Println("URL file name must be specified")

		return
	}

	if len(*fieldNameFlag) == 0 {
		fmt.Println("Field name file name must be specified")

		return
	}

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	parsedData, err := url.ParseQuery(*dataFlag)

	if err != nil {
		fmt.Println("Error parsing data:", err)

		return
	}

	for k, v := range parsedData {
		if err = writer.WriteField(k, v[0]); err != nil {
			fmt.Println("Error writing data:", err)

			return
		}
	}

	for _, filePath := range filePaths {
		err := addFileToWriter(writer, filePath, *fieldNameFlag)

		if err != nil {
			fmt.Println("Error adding file:", err)

			return
		}
	}

	err = writer.Close()

	if err != nil {
		fmt.Println("Error closing writer:", err)

		return
	}

	resp, err := http.Post(*urlFlag, writer.FormDataContentType(), &body)

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
