package agent

import (
	"context"
	"encoding/json"
	"example.com/internal/config"
	"example.com/internal/models"
	"fmt"
	"github.com/google/uuid"
	"log"
	"runtime/debug"
	"time"
)

func GetPersonalInfoInJsonFormat(resumeDetails string) (*models.PersonalInfo, error) {
	client := GetClientForResumeParserAgent()
	log.Printf("Model used for parsing Personal Info %v", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetPersonalInfoSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": models.PersonalInfo{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing personal info: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for personal info")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var personalInfo models.PersonalInfo
	err = json.Unmarshal([]byte(contentStr), &personalInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling personal info: %v", err)
	}

	return &personalInfo, nil
}

func GetEducationInJsonFormat(resumeDetails string) ([]models.Education, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing Education %v", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetEducationSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": struct {
					Education []models.Education `json:"education"`
				}{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing education: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for education")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var result struct {
		Education []models.Education `json:"education"`
	}
	err = json.Unmarshal([]byte(contentStr), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling education: %v", err)
	}

	return result.Education, nil
}

func GetExperienceInJsonFormat(resumeDetails string) ([]models.Experience, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing Experience %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetExperienceSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": struct {
					Experience []models.Experience `json:"experience"`
				}{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing experience: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for experience")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var result struct {
		Experience []models.Experience `json:"experience"`
	}
	err = json.Unmarshal([]byte(contentStr), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling experience: %v", err)
	}

	return result.Experience, nil
}

func GetSkillsInJsonFormat(resumeDetails string) (*models.Skills, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing skills %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetSkillsSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": models.Skills{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing skills: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for skills")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var skills models.Skills
	err = json.Unmarshal([]byte(contentStr), &skills)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling skills: %v", err)
	}

	return &skills, nil
}

func GetProjectsInJsonFormat(resumeDetails string) ([]models.Project, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing projects %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetProjectsSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": struct {
					Projects []models.Project `json:"projects"`
				}{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing projects: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for projects")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var result struct {
		Projects []models.Project `json:"projects"`
	}
	err = json.Unmarshal([]byte(contentStr), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling projects: %v", err)
	}

	return result.Projects, nil
}

func GetSummaryInJsonFormat(resumeDetails string) (string, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing summary %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetSummarySystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": struct {
					Summary string `json:"summary"`
				}{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error parsing summary: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response received for summary")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return "", fmt.Errorf("response content is not a string")
	}

	var result struct {
		Summary string `json:"summary"`
	}
	err = json.Unmarshal([]byte(contentStr), &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling summary: %v", err)
	}

	return result.Summary, nil
}

func GetAdditionalInfoInJsonFormat(resumeDetails string) (*models.AdditionalInfo, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing additional info %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetAdditionalInfoSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": models.AdditionalInfo{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing additional info: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for additional info")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var additionalInfo models.AdditionalInfo
	err = json.Unmarshal([]byte(contentStr), &additionalInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling additional info: %v", err)
	}

	return &additionalInfo, nil
}

func GetHobbiesInJsonFormat(resumeDetails string) ([]string, error) {
	client := GetClientForResumeParserAgent()

	log.Printf("Model used for parsing hobbies %v", client.config.Model)
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetHobbiesSystemPrompt(),
			},
			{
				Role:    "user",
				Content: resumeDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": struct {
					Hobbies []string `json:"hobbies"`
				}{},
			},
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error parsing hobbies: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received for hobbies")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	var result struct {
		Hobbies []string `json:"hobbies"`
	}
	err = json.Unmarshal([]byte(contentStr), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling hobbies: %v", err)
	}

	return result.Hobbies, nil
}

// Main function to assemble the complete Resume
func GetResumeDetailInJsonFormat(resumeDetails string) (*models.Resume, error) {
	title := "Sample Resume Title"
	description := "Sample Resume Description"
	fileType := "application/pdf"
	// Parse each section concurrently for better performance
	type result struct {
		personalInfo   *models.PersonalInfo
		education      []models.Education
		experience     []models.Experience
		skills         *models.Skills
		projects       []models.Project
		summary        string
		additionalInfo *models.AdditionalInfo
		hobbies        []string
		err            error
	}

	ch := make(chan result, 8)

	// Launch goroutines for each parsing task
	go func() {
		personalInfo, err := GetPersonalInfoInJsonFormat(resumeDetails)
		ch <- result{personalInfo: personalInfo, err: err}
	}()

	go func() {
		education, err := GetEducationInJsonFormat(resumeDetails)
		ch <- result{education: education, err: err}
	}()

	go func() {
		experience, err := GetExperienceInJsonFormat(resumeDetails)
		ch <- result{experience: experience, err: err}
	}()

	go func() {
		skills, err := GetSkillsInJsonFormat(resumeDetails)
		ch <- result{skills: skills, err: err}
	}()

	go func() {
		projects, err := GetProjectsInJsonFormat(resumeDetails)
		ch <- result{projects: projects, err: err}
	}()

	go func() {
		summary, err := GetSummaryInJsonFormat(resumeDetails)
		ch <- result{summary: summary, err: err}
	}()

	go func() {
		additionalInfo, err := GetAdditionalInfoInJsonFormat(resumeDetails)
		ch <- result{additionalInfo: additionalInfo, err: err}
	}()

	go func() {
		hobbies, err := GetHobbiesInJsonFormat(resumeDetails)
		ch <- result{hobbies: hobbies, err: err}
	}()

	// Collect results
	var (
		personalInfo   *models.PersonalInfo
		education      []models.Education
		experience     []models.Experience
		skills         *models.Skills
		projects       []models.Project
		summary        string
		additionalInfo *models.AdditionalInfo
		hobbies        []string
		errors         []error
	)

	for i := 0; i < 8; i++ {
		res := <-ch
		if res.err != nil {
			errors = append(errors, res.err)
			continue
		}

		if res.personalInfo != nil {
			personalInfo = res.personalInfo
		}
		if res.education != nil {
			education = res.education
		}
		if res.experience != nil {
			experience = res.experience
		}
		if res.skills != nil {
			skills = res.skills
		}
		if res.projects != nil {
			projects = res.projects
		}
		if res.summary != "" {
			summary = res.summary
		}
		if res.additionalInfo != nil {
			additionalInfo = res.additionalInfo
		}
		if res.hobbies != nil {
			hobbies = res.hobbies
		}
	}

	// Log errors but continue with partial data
	if len(errors) > 0 {
		for _, err := range errors {
			log.Printf("Parsing error: %v", err)
		}
	}

	// Create the complete Resume object
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		newUUID = uuid.Nil
	}
	resume := &models.Resume{
		UserID:      newUUID,
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

	// Set embedded structs
	if personalInfo != nil {
		resume.PersonalInfo = *personalInfo
	}
	if additionalInfo != nil {
		resume.AdditionalInfo = additionalInfo
	}

	return resume, nil
}

// Alternative sequential approach (if you prefer not to use goroutines)
func GetResumeDetailInJsonFormatSequential(resumeDetails string, userID uuid.UUID, title, description, fileType string) (*models.Resume, error) {
	resume := &models.Resume{
		UserID:      userID,
		Title:       title,
		Description: description,
		FileType:    fileType,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Parse personal info
	if personalInfo, err := GetPersonalInfoInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing personal info: %v", err)
	} else if personalInfo != nil {
		log.Printf("personal info parsing sucessfull")
		resume.PersonalInfo = *personalInfo
	}

	// Parse education
	if education, err := GetEducationInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing education: %v", err)
	} else {
		log.Printf("education parsing sucessfull")
		resume.Education = education
	}

	// Parse experience
	if experience, err := GetExperienceInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing experience: %v", err)
	} else {
		log.Printf("experience parsing sucessfull")
		resume.Experience = experience
	}

	// Parse skills
	if skills, err := GetSkillsInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing skills: %v", err)
	} else {
		log.Printf("skills parsing sucessfull")
		resume.Skills = skills
	}

	// Parse projects
	if projects, err := GetProjectsInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing projects: %v", err)
	} else {
		log.Printf("projects parsing sucessfull")
		resume.Projects = projects
	}

	// Parse summary
	if summary, err := GetSummaryInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing summary: %v", err)
	} else {
		log.Printf("summary parsing sucessfull")
		resume.Summary = summary
	}

	// Parse additional info
	if additionalInfo, err := GetAdditionalInfoInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing additional info: %v", err)
	} else {
		log.Printf("additional info parsing sucessfull")
		resume.AdditionalInfo = additionalInfo
	}

	// Parse hobbies
	if hobbies, err := GetHobbiesInJsonFormat(resumeDetails); err != nil {
		log.Printf("Error parsing hobbies: %v", err)
	} else {
		log.Printf("hobbies parsing sucessfull")
		resume.Hobbies = hobbies
	}

	return resume, nil
}

// System prompt functions with detailed schema specifications
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
func GetClientForResumeParserAgent() *ConfigAgentClient {
	cfg, err := config.LoadAgentClient()

	if err != nil {
		log.Printf("Failed to load agent client configuration: %s\nStackTrace:\n%s", err.Error(), debug.Stack())
		return nil
	}

	client := NewAgentClient(cfg)
	return client
}

func ResumeFixJsonFormat(resumeDetails string, errorDetails string) string {
	client := GetClientForResumeParserAgent()

	req := &ChatRequest{ // Use pointer to ChatRequest
		Messages: []Message{
			{
				Role:    "system",
				Content: models.GetResumeJsonFormatFix(),
			},
			{
				Role:    "user",
				Content: resumeDetails + "\n" + errorDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
			Schema: map[string]interface{}{
				"schema": models.Resume{},
			},
		},
	}
	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print response
	if len(response.Choices) > 0 {
		fmt.Printf("AI Response: %v\n", response.Choices[0].Message.Content)
		return fmt.Sprintf("%v", response.Choices[0].Message.Content)
	} else {
		fmt.Println("No response received")
		return ""
	}
}
