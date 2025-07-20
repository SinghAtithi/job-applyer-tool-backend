package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"example.com/agent"
	"example.com/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
)

const (
	MaxFileSize       = 10 << 20 // 10 MB
	TextPreviewLength = 200
	SupportedFileExt  = ".pdf"
	MaxRetryAttempts  = 3 // Maximum attempts to fix JSON format
)

// PDFExtractor handles PDF text extraction
type PDFExtractor struct{}

// ExtractText extracts text content from a PDF byte slice
func (pe *PDFExtractor) ExtractText(fileContent []byte) (string, error) {
	reader := bytes.NewReader(fileContent)

	pdfReader, err := pdf.NewReader(reader, int64(len(fileContent)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var textBuilder strings.Builder

	for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
		if err := pe.extractPageText(pdfReader, pageNum, &textBuilder); err != nil {
			log.Printf("Warning: failed to extract text from page %d: %v", pageNum, err)
			continue
		}
	}

	return textBuilder.String(), nil
}

// extractPageText extracts text from a single page
func (pe *PDFExtractor) extractPageText(pdfReader *pdf.Reader, pageNum int, textBuilder *strings.Builder) error {
	page := pdfReader.Page(pageNum)
	if page.V.IsNull() {
		return nil // Skip null pages
	}

	textContent, err := page.GetPlainText(nil)
	if err != nil {
		return fmt.Errorf("failed to get plain text: %w", err)
	}

	textBuilder.WriteString(textContent)
	textBuilder.WriteString("\n")
	return nil
}

// FileProcessor handles file processing operations
type FileProcessor struct {
	extractor *PDFExtractor
}

// NewFileProcessor creates a new file processor
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		extractor: &PDFExtractor{},
	}
}

// ProcessFile processes an uploaded file and extracts text content
func (fp *FileProcessor) ProcessFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

	if !fp.isSupportedFileType(header.Filename) {
		return "", fmt.Errorf("unsupported file type: %s (only %s files are supported)", header.Filename, SupportedFileExt)
	}

	fileContent, err := fp.readFileContent(file, header.Size)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	log.Printf("Processing file: %s (Size: %d bytes)", header.Filename, header.Size)

	textContent, err := fp.extractor.ExtractText(fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	fp.logExtractedText(textContent)
	return textContent, nil
}

// isSupportedFileType checks if the file type is supported
func (fp *FileProcessor) isSupportedFileType(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), SupportedFileExt)
}

// readFileContent reads the entire file content into memory
func (fp *FileProcessor) readFileContent(file multipart.File, size int64) ([]byte, error) {
	content := make([]byte, size)
	_, err := io.ReadFull(file, content)
	return content, err
}

// logExtractedText logs a preview of the extracted text
func (fp *FileProcessor) logExtractedText(textContent string) {
	if len(textContent) > TextPreviewLength {
		log.Printf("Extracted text preview: %s...", textContent[:TextPreviewLength])
	} else {
		log.Printf("Extracted text: %s", textContent)
	}
}

// ResumeParser handles resume parsing operations
type ResumeParser struct {
	fileProcessor *FileProcessor
}

// NewResumeParser creates a new resume parser
func NewResumeParser() *ResumeParser {
	return &ResumeParser{
		fileProcessor: NewFileProcessor(),
	}
}

// ParseFromUpload processes an uploaded file and converts it to a Resume struct
func (rp *ResumeParser) ParseFromUpload(c *gin.Context) (models.Resume, error) {
	if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		return models.Resume{}, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return models.Resume{}, fmt.Errorf("failed to retrieve file from form: %w", err)
	}

	textContent, err := rp.fileProcessor.ProcessFile(file, header)
	if err != nil {
		return models.Resume{}, err
	}

	return rp.parseTextToResume(textContent)
}

// parseTextToResume converts extracted text to a Resume struct
func (rp *ResumeParser) parseTextToResume(textContent string) (models.Resume, error) {
	resumeDetails, err := agent.GetResumeDetailInJsonFormat(textContent)
	if err != nil {
		return models.Resume{}, fmt.Errorf("failed to parse resume details: %w", err)
	}

	log.Printf("Successfully parsed resume: %+v", resumeDetails)
	return *resumeDetails, nil
}

// parseJSONToResume converts a JSON string to a Resume struct
func (rp *ResumeParser) parseJSONToResume(jsonStr string) (*models.Resume, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return nil, fmt.Errorf("JSON string is empty")
	}

	jsonStr = strings.TrimSpace(jsonStr)
	if !rp.isValidJSONFormat(jsonStr) {
		return nil, fmt.Errorf("invalid JSON format: must be a valid JSON object")
	}

	// Validate JSON syntax first
	if err := rp.validateJSONSyntax(jsonStr); err != nil {
		return nil, err
	}

	var resume models.Resume
	if err := json.Unmarshal([]byte(jsonStr), &resume); err != nil {
		return nil, rp.createUnmarshalError(err)
	}

	return &resume, nil
}

// isValidJSONFormat performs basic JSON format validation
func (rp *ResumeParser) isValidJSONFormat(jsonStr string) bool {
	return strings.HasPrefix(jsonStr, "{") && strings.HasSuffix(jsonStr, "}")
}

// validateJSONSyntax validates JSON syntax without unmarshaling to specific struct
func (rp *ResumeParser) validateJSONSyntax(jsonStr string) error {
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("JSON syntax error at position %d: %w", syntaxErr.Offset, err)
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// createUnmarshalError creates detailed error messages for unmarshal failures
func (rp *ResumeParser) createUnmarshalError(err error) error {
	switch e := err.(type) {
	case *json.SyntaxError:
		return fmt.Errorf("JSON syntax error at position %d: %w", e.Offset, err)
	case *json.UnmarshalTypeError:
		return fmt.Errorf("type mismatch: field '%s' expects %s but got %s at position %d",
			e.Field, e.Type, e.Value, e.Offset)
	case *json.UnsupportedTypeError:
		return fmt.Errorf("unsupported type: %s", e.Type)
	default:
		return fmt.Errorf("failed to unmarshal JSON to Resume struct: %w", err)
	}
}

// Public API functions for backward compatibility and ease of use

// HandleCreateResume processes file upload and creates a Resume struct
func HandleCreateResume(c *gin.Context) (models.Resume, error) {
	parser := NewResumeParser()
	return parser.ParseFromUpload(c)
}

// GetResumeSchemaFromJsonTyped converts JSON string to Resume struct
func GetResumeSchemaFromJsonTyped(jsonStr string) (*models.Resume, error) {
	parser := NewResumeParser()
	return parser.parseJSONToResume(jsonStr)
}

func (rp *ResumeParser) ResumeToText(ctx *gin.Context) (string, error) {
	if err := ctx.Request.ParseMultipartForm(MaxFileSize); err != nil {
		return "", fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file from form: %w", err)
	}

	textContent, err := rp.fileProcessor.ProcessFile(file, header)
	if err != nil {
		return "", err
	}
	return textContent, nil
}

// HandleParseResume placeholder for future implementation
func HandleParseResume(ctx *gin.Context) (models.Resume, error) {
	resumeParser := NewResumeParser()
	resumeString, err := resumeParser.ResumeToText(ctx)
	if err != nil {
		return models.Resume{}, fmt.Errorf("failed to parse resume: %w", err)
	}

	resumeDetails, parseErr := agent.GetResumeDetailInJsonFormat(resumeString)
	resumeInJsonFormat := ""
	if parseErr != nil {
		resumeInJsonFormat = agent.ResumeFixJsonFormat(resumeString, parseErr.Error())
	} else if resumeDetails != nil {
		resumeInJsonFormatBytes, err := json.Marshal(resumeDetails)
		if err != nil {
			return models.Resume{}, fmt.Errorf("failed to marshal resume details: %w", err)
		}
		resumeInJsonFormat = string(resumeInJsonFormatBytes)
	}

	err = resumeParser.validateJSONSyntax(resumeInJsonFormat)

	if err != nil {
		attempts := 0
		maxAttempts := MaxRetryAttempts
		for attempts < maxAttempts {
			resumeInJsonFormat = agent.ResumeFixJsonFormat(resumeString, err.Error())
			attempts++
			err = resumeParser.validateJSONSyntax(resumeInJsonFormat)
			if err == nil {
				break
			}
		}
	}

	resumeDetailsObj, err := GetResumeSchemaFromJsonTyped(resumeInJsonFormat)

	if err != nil {
		attempts := 0
		maxAttempts := MaxRetryAttempts
		for attempts < maxAttempts {
			resumeInJsonFormat = agent.ResumeFixJsonFormat(resumeString, err.Error())
			attempts++
			resumeDetailsObj, err = GetResumeSchemaFromJsonTyped(resumeInJsonFormat)
			if err == nil {
				break
			}
		}
	}
	if err != nil || resumeDetailsObj == nil {
		return models.Resume{}, errors.New("failed to parse resume details after multiple attempts")
	}
	return *resumeDetailsObj, err
}
