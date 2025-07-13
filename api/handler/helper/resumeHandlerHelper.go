package helper

import (
	"bytes"
	"example.com/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"log"
	"mime/multipart"
	_ "net/http"
	"strings"
)

func extractTextFromPDF(fileContent []byte) (string, error) {
	// Create a reader from the byte slice
	reader := bytes.NewReader(fileContent)

	// Open the PDF
	pdfReader, err := pdf.NewReader(reader, int64(len(fileContent)))
	if err != nil {
		return "", fmt.Errorf("error creating PDF reader: %w", err)
	}

	var textBuilder strings.Builder

	// Extract text from each page
	for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Get text content from the page
		textContent, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("Error extracting text from page %d: %v", pageNum, err)
			continue
		}

		textBuilder.WriteString(textContent)
		textBuilder.WriteString("\n") // Add newline between pages
	}

	return textBuilder.String(), nil
}

func HandleCreateResume(c *gin.Context) (models.Resume, error) {
	// Parse the multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return models.Resume{}, fmt.Errorf("error parsing form: %w", err)
	}

	// Get the file from the form
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		return models.Resume{}, fmt.Errorf("error retrieving file: %w", err)
	}

	// Extract the file content to a variable
	fileContentBinary := make([]byte, handler.Size)

	// Read the file content into the variable
	_, err = file.Read(fileContentBinary)
	if err != nil {
		return models.Resume{}, fmt.Errorf("error reading file content: %w", err)
	}

	// Log the file name and size
	log.Printf("Received file: %s, Size: %d bytes", handler.Filename, handler.Size)

	// Ensure the file is closed after processing
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Extract text content from PDF
	var textContent string
	if strings.HasSuffix(strings.ToLower(handler.Filename), ".pdf") {
		textContent, err = extractTextFromPDF(fileContentBinary)
		if err != nil {
			return models.Resume{}, fmt.Errorf("error extracting text from PDF: %w", err)
		}
	} else {
		return models.Resume{}, fmt.Errorf("unsupported file type: %s", handler.Filename)
	}

	// Get the other form values
	//userID := c.Request.FormValue("userID")
	//title := c.Request.FormValue("title")
	//description := c.Request.FormValue("description")

	// Log the extracted text (first 200 characters for brevity)
	if len(textContent) > 200 {
		log.Printf("Extracted text preview: %s...", textContent[:200])
	} else {
		log.Printf("Extracted text: %s", textContent)
	}

	return models.Resume{
		PersonalInfo: models.PersonalInfo{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "Jhon@Doe.com",
			Phone:     "123-456-7890",
			LinkedIn:  "https://www.linkedin.com/in/johndoe",
			GitHub:    "https://www.github.com/johndoe",
			Website:   "https://www.johndoe.com",
		},
	}, nil
}
