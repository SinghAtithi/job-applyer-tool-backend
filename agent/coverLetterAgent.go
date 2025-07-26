package agent

import (
	"context"
	"fmt"
	"log"
)

func GetJobDescription(jobDescription string) (string, error) {
	client := GetClientForResumeParserAgent()
	log.Printf("Model used for parsing Job Description Info %v", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: getJobDescriptionContentPrompt(),
			},
			{
				Role:    "user",
				Content: jobDescription,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "text",
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error parsing personal info: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response received for personal info")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return "", fmt.Errorf("response content is not a string")
	}
	return contentStr, nil
}

func GetCoverLetterString(jobDescription string, userDetails string) (string, error) {
	client := GetClientForResumeParserAgent()
	log.Printf("Model used for parsing Cover Letter Info %v", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: getCoverLetterPrompt(),
			},
			{
				Role:    "user",
				Content: jobDescription + "\n" + userDetails,
			},
		},
		ResponseFormat: &ResponseFormat{
			Type: "text",
		},
	}

	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error parsing personal info: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response received for personal info")
	}

	contentStr, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return "", fmt.Errorf("response content is not a string")
	}
	return contentStr, nil
}

func getJobDescriptionContentPrompt() string {
	return `
Extract the following information from the job description text:
1. Extract basic job details like title, company, location, employment type
2. Identify educational qualifications, experience requirements, and certifications needed
3. List technical skills, programming languages, and software tools required
4. Extract soft skills and personal attributes mentioned
5. Summarize key job responsibilities and duties

## Instructions:
- Extract only factual information present in the text
- Ignore promotional content and irrelevant details  
- If information is not available, mark as "Not specified"
- Present in clean, organized format
- Output only the extracted data with No additional text, explanations, or commentary
- No greetings like "Here's your job description" or similar phrases
- Start directly with the Job Description
- The response will be used directly for further operation`
}

func getCoverLetterPrompt() string {
	return `You will receive a job description and a JSON object containing applicant details. Generate a professional cover letter based on this information.
Requirements:
Write a formal cover letter that is exactly 3-4 paragraphs long
Include the applicant's relevant skills and experience from the JSON data
Align the applicant's qualifications with the job requirements from the description
Explain how the applicant's background can benefit the specific role
Maintain a professional and confident tone throughout
Critical Instructions:
Output ONLY the cover letter content
No additional text, explanations, or commentary
No greetings like "Here's your cover letter" or similar phrases
Start directly with the cover letter opening
End with the cover letter closing
The response will be used directly for PDF generation
Format:
The cover letter should follow standard business format with proper paragraph structure and professional language.
`
}
