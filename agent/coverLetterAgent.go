package agent

import (
	"context"
	"fmt"
	"log"
)

func GetCoverLetterContent(jobDescription string) (string, error) {
	client := GetClientForResumeParserAgent()
	log.Printf("Model used for parsing Personal Info %v", client.config.Model)

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: GetCoverLetterContentPrompt(),
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

func GetCoverLetterContentPrompt() string {
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
- Output only the extracted data with no extra words, explanations, or commentary`
}
