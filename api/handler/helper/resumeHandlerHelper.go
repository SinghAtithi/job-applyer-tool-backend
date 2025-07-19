package helper

import (
	"bytes"
	"encoding/json"
	"example.com/agent"
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

	parsedResumeDetails := agent.GetResumeDetailInJsonFormat(textContent)

	// parsedResumeDetails has the resume details in String format, first convert it to Resume struct
	resumeDetails, err := GetResumeSchemaFromJsonTyped(parsedResumeDetails)
	if resumeDetails == nil || err != nil {
		return models.Resume{}, fmt.Errorf("error parsing resume details from JSON string: %w", err)
	}

	log.Printf("Resume details: %+v", resumeDetails)
	return *resumeDetails, nil

}

func GetResumeSchemaFromJsonTyped(details string) (*models.Resume, error) {
	// Input validation
	if details == "" {
		return nil, fmt.Errorf("input JSON string is empty")
	}

	// Trim whitespace and check for basic JSON structure
	details = strings.TrimSpace(details)
	if !strings.HasPrefix(details, "{") || !strings.HasSuffix(details, "}") {
		return nil, fmt.Errorf("invalid JSON format: must start with '{' and end with '}'")
	}

	// Check if it's valid JSON first
	var jsonCheck interface{}
	if err := json.Unmarshal([]byte(details), &jsonCheck); err != nil {
		// More specific JSON syntax error
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("JSON syntax error at position %d: %w", syntaxErr.Offset, err)
		}
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("JSON type error: cannot unmarshal %s into field %s of type %s at position %d",
				typeErr.Value, typeErr.Field, typeErr.Type, typeErr.Offset)
		}
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Initialize the resume struct
	resumeDetails := &models.Resume{}

	// Attempt to unmarshal into the specific struct
	err := json.Unmarshal([]byte(details), resumeDetails)
	if err != nil {
		// Provide detailed error information
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("JSON syntax error while parsing resume at position %d: %w", syntaxErr.Offset, err)
		}

		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("type mismatch in resume JSON: field '%s' expects type %s but got %s at position %d",
				typeErr.Field, typeErr.Type, typeErr.Value, typeErr.Offset)
		}

		if fieldErr, ok := err.(*json.UnsupportedTypeError); ok {
			return nil, fmt.Errorf("unsupported type error in resume JSON: %s", fieldErr.Type)
		}

		return nil, fmt.Errorf("failed to parse resume JSON into struct: %w", err)
	}

	return resumeDetails, nil
}
