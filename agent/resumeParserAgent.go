package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example.com/internal/models"
	"example.com/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func parseSection[T any](ctx context.Context, systemPrompt, userInput string, schema map[string]interface{}) (T, error) {
	var zero T
	client, err := GetSharedClient()
	if err != nil {
		return zero, fmt.Errorf("failed to get agent client: %w", err)
	}

	logger.Info("parsing section (model: %s)", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userInput},
		},
		ResponseFormat: &ResponseFormat{
			Type:   "json_object",
			Schema: schema,
		},
	}

	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("LLM request failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return zero, fmt.Errorf("no response received")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return zero, fmt.Errorf("response content is not a string")
	}

	var result T
	if err := json.Unmarshal([]byte(contentStr), &result); err != nil {
		return zero, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return result, nil
}

func GetPersonalInfoInJsonFormat(ctx context.Context, resumeDetails string) (*models.PersonalInfo, error) {
	return parseSection[*models.PersonalInfo](ctx, GetPersonalInfoSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": models.PersonalInfo{}})
}

func GetEducationInJsonFormat(ctx context.Context, resumeDetails string) ([]models.Education, error) {
	type wrapper struct {
		Education []models.Education `json:"education"`
	}
	result, err := parseSection[wrapper](ctx, GetEducationSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": wrapper{}})
	if err != nil {
		return nil, err
	}
	return result.Education, nil
}

func GetExperienceInJsonFormat(ctx context.Context, resumeDetails string) ([]models.Experience, error) {
	type wrapper struct {
		Experience []models.Experience `json:"experience"`
	}
	result, err := parseSection[wrapper](ctx, GetExperienceSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": wrapper{}})
	if err != nil {
		return nil, err
	}
	return result.Experience, nil
}

func GetSkillsInJsonFormat(ctx context.Context, resumeDetails string) (*models.Skills, error) {
	return parseSection[*models.Skills](ctx, GetSkillsSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": models.Skills{}})
}

func GetProjectsInJsonFormat(ctx context.Context, resumeDetails string) ([]models.Project, error) {
	type wrapper struct {
		Projects []models.Project `json:"projects"`
	}
	result, err := parseSection[wrapper](ctx, GetProjectsSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": wrapper{}})
	if err != nil {
		return nil, err
	}
	return result.Projects, nil
}

func GetSummaryInJsonFormat(ctx context.Context, resumeDetails string) (string, error) {
	type wrapper struct {
		Summary string `json:"summary"`
	}
	result, err := parseSection[wrapper](ctx, GetSummarySystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": wrapper{}})
	if err != nil {
		return "", err
	}
	return result.Summary, nil
}

func GetAdditionalInfoInJsonFormat(ctx context.Context, resumeDetails string) (*models.AdditionalInfo, error) {
	return parseSection[*models.AdditionalInfo](ctx, GetAdditionalInfoSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": models.AdditionalInfo{}})
}

func GetHobbiesInJsonFormat(ctx context.Context, resumeDetails string) ([]string, error) {
	type wrapper struct {
		Hobbies []string `json:"hobbies"`
	}
	result, err := parseSection[wrapper](ctx, GetHobbiesSystemPrompt(), resumeDetails,
		map[string]interface{}{"schema": wrapper{}})
	if err != nil {
		return nil, err
	}
	return result.Hobbies, nil
}

func GetResumeDetailInJsonFormat(ctx context.Context, resumeDetails string, userID uuid.UUID, title, description, fileType string) (*models.Resume, error) {
	g, ctx := errgroup.WithContext(ctx)

	var (
		personalInfo   *models.PersonalInfo
		education      []models.Education
		experience     []models.Experience
		skills         *models.Skills
		projects       []models.Project
		summary        string
		additionalInfo *models.AdditionalInfo
		hobbies        []string
	)

	g.Go(func() error {
		var err error
		personalInfo, err = GetPersonalInfoInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse personal info: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		education, err = GetEducationInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse education: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		experience, err = GetExperienceInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse experience: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		skills, err = GetSkillsInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse skills: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		projects, err = GetProjectsInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse projects: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		summary, err = GetSummaryInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse summary: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		additionalInfo, err = GetAdditionalInfoInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse additional info: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		hobbies, err = GetHobbiesInJsonFormat(ctx, resumeDetails)
		if err != nil {
			logger.Warn("failed to parse hobbies: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("resume parsing failed: %w", err)
	}

	resume := &models.Resume{
		UserID:      userID,
		Title:       title,
		Description: description,
		FileType:    fileType,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Summary:     summary,
		Education:   education,
		Experience:  experience,
		Skills:      skills,
		Projects:    projects,
		Hobbies:     hobbies,
	}

	if personalInfo != nil {
		resume.PersonalInfo = *personalInfo
	}
	if additionalInfo != nil {
		resume.AdditionalInfo = additionalInfo
	}

	return resume, nil
}

func ResumeFixJsonFormat(ctx context.Context, resumeDetails string, errorDetails string) (string, error) {
	client, err := GetSharedClient()
	if err != nil {
		return "", fmt.Errorf("failed to get agent client: %w", err)
	}

	req := &ChatRequest{
		Messages: []Message{
			{Role: "system", Content: models.GetResumeJsonFormatFix()},
			{Role: "user", Content: resumeDetails + "\n" + errorDetails},
		},
		ResponseFormat: &ResponseFormat{
			Type:   "json_object",
			Schema: map[string]interface{}{"schema": models.Resume{}},
		},
	}

	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to fix JSON format via AI: %w", err)
	}

	if len(response.Choices) > 0 {
		contentStr, ok := response.Choices[0].Message.Content.(string)
		if ok {
			return contentStr, nil
		}
		return fmt.Sprintf("%v", response.Choices[0].Message.Content), nil
	}

	return "", fmt.Errorf("no response received for JSON format fix")
}

func GetPersonalInfoSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract personal information from the provided resume text and return ONLY a valid JSON object that matches the PersonalInfo schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found in the resume, omit it from the JSON (do not use null or empty strings)
4. All date fields must be in YYYY-MM-DDTHH:MM:SSZ format (ISO 8601 with Z timezone)
5. Return ONLY the JSON object without any code blocks or markdown formatting
6. If date_of_birth is not provided, omit it from the JSON
7. Middle Name if present should be included with first_name and not with last_name

PersonalInfo Schema:
{
  "first_name": "string",
  "last_name": "string", 
  "middle_name": "string",
  "title": "string",
  "email": "string",
  "phone": "string",
  "address": {
    "street": "string",
    "city": "string", 
    "state": "string",
    "country": "string",
    "postal_code": "string"
  },
  "website": "string",
  "linkedin": "string", 
  "github": "string",
  "twitter": "string",
  "portfolio": "string",
  "date_of_birth": "YYYY-MM-DDTHH:MM:SSZ",
  "nationality": "string",
  "gender": "string",
  "marital_status": "string"
}

Extract personal information and return ONLY the JSON object.`
}

func GetEducationSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract education information from the provided resume text and return ONLY a valid JSON object containing an education array that matches the Education schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found, omit it from the JSON (do not use null or empty strings)
4. All date fields must be in YYYY-MM-DDTHH:MM:SSZ format (ISO 8601 with Z timezone)
5. Return the result as {"education": [...]} format
6. Return ONLY the JSON object without any code blocks or markdown formatting
7. Id should be a unique integer starting from 1 and incrementing for each education entry

Education Schema (array of objects):
{
  "education": [
    {
      "id": "number",
      "institution": "string",
      "degree": "string",
      "field_of_study": "string", 
      "grade": "string",
      "grade_type": "string",
      "start_date": "YYYY-MM-DDTHH:MM:SSZ",
      "end_date": "YYYY-MM-DDTHH:MM:SSZ",
      "is_current": "boolean"
    }
  ]
}

Extract education information and return ONLY the JSON object.`
}

func GetExperienceSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract work experience from the provided resume text and return ONLY a valid JSON object containing an experience array that matches the Experience schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found, omit it from the JSON (do not use null or empty strings)
4. All date fields must be in YYYY-MM-DDTHH:MM:SSZ format (ISO 8601 with Z timezone)
5. Return the result as {"experience": [...]} format
6. Return ONLY the JSON object without any code blocks or markdown formatting
7. technicalSkill.technical.experience should be of type int


Experience Schema (array of objects):
{
  "experience": [
    {
      "id": "number",
      "job_title": "string",
      "company": "string",
      "role": "string",
      "location": "string",
      "start_date": "YYYY-MM-DDTHH:MM:SSZ",
      "end_date": "YYYY-MM-DDTHH:MM:SSZ", 
      "is_current": "boolean",
      "employment_type": "string",
      "description": "string",
      "responsibilities": ["string", "string"],
      "achievements": ["string", "string"],
      "technologies_used": ["string", "string"]
    }
  ]
}

Extract work experience and return ONLY the JSON object.`
}

func GetSkillsSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract skills information from the provided resume text and return ONLY a valid JSON object that matches the Skills schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found, omit it from the JSON (do not use null or empty strings)
4. Categorize skills appropriately into the provided categories
5. Return ONLY the JSON object without any code blocks or markdown formatting

Skills Schema:
{
  "technical": [
    {
      "name": "string",
      "level": "string",
      "experience": "number",
      "category": "string"
    }
  ],
  "tools": ["string", "string"],
  "frameworks": ["string", "string"],
  "libraries": ["string", "string"],
  "databases": ["string", "string"],
  "cloud_platforms": ["string", "string"],
  "operating_systems": ["string", "string"],
  "methodologies": ["string", "string"]
}

Extract skills information and return ONLY the JSON object.`
}

func GetProjectsSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract project information from the provided resume text and return ONLY a valid JSON object containing a projects array that matches the Project schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found, omit it from the JSON (do not use null or empty strings)
4. All date fields must be in YYYY-MM-DDTHH:MM:SSZ format (ISO 8601 with Z timezone)
5. Return the result as {"projects": [...]} format
6. Return ONLY the JSON object without any code blocks or markdown formatting
7. If start_date or end_date is not provided, set them to null or omit them

Project Schema (array of objects):
{
  "projects": [
    {
      "id": "number",
      "name": "string",
      "description": "string",
      "technologies": ["string", "string"],
      "url": "string",
      "github_url": "string",
      "start_date": "YYYY-MM-DDTHH:MM:SSZ",
      "end_date": "YYYY-MM-DDTHH:MM:SSZ",
      "is_current": "boolean"
    }
  ]
}

Extract project information and return ONLY the JSON object.`
}

func GetSummarySystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract or identify the professional summary from the provided resume text and return ONLY a valid JSON object. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Extract existing summary/objective from resume text - do not create or generate new content
3. If no summary is found, return {"summary": ""}
4. Return ONLY the JSON object without any code blocks or markdown formatting

Summary Schema:
{
  "summary": "string"
}

Extract the professional summary and return ONLY the JSON object.`
}

func GetAdditionalInfoSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract additional information from the provided resume text and return ONLY a valid JSON object that matches the AdditionalInfo schema. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. If a field is not found, omit it from the JSON (do not use null or empty strings)
4. Return ONLY the JSON object without any code blocks or markdown formatting
5. Field types should be strictly followed: do not convert booleans to strings or vice versa and same with integers

AdditionalInfo Schema:
{
  "availability": "string",
  "notice_period": "string",
  "salary": "string",
  "willing_to_relocate": "boolean",
  "willing_to_travel": "string",
  "work_authorization": "string",
  "security_clearance": "string",
  "driving_license": "string", 
  "military_service": "string",
  "disabilities": "string",
  "emergency_contact": "string",
  "preferred_location": "string",
  "remote_work": "boolean",
  "custom_fields": {}
}

Extract additional information and return ONLY the JSON object.`
}

func GetHobbiesSystemPrompt() string {
	return `You are a resume data extraction system. Your task is to extract hobbies and interests from the provided resume text and return ONLY a valid JSON object containing a hobbies array. Follow these rules strictly:

1. Return ONLY valid JSON - no additional text, explanations, or formatting
2. Use only the data present in the resume text - do not create, assume, or invent any information
3. Look for sections like "Interests", "Hobbies", "Personal Interests", "Activities", etc.
4. If no hobbies are found, return {"hobbies": []}
5. Return ONLY the JSON object without any code blocks or markdown formatting

Hobbies Schema:
{
  "hobbies": ["string", "string", "string"]
}

Extract hobbies and interests and return ONLY the JSON object.`
}