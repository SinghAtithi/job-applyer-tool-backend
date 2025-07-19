package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"example.com/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func MakeAgentCall(prompt string) string {
	// Get the API key from environment variable
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Println("GROQ_API_KEY not set in environment")
		return "API key not set"
	}
	url := "https://api.groq.com/openai/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "llama3-8b-8192",
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "You are a resume data extraction system. Your task is to parse the provided resume text and return ONLY a JSON object that matches the following schema. Follow these rules strictly:\n\n1. Return ONLY valid JSON - no additional text, explanations, or formatting\n2. Use only the data present in the provided text - do not create, assume, or invent any information\n3. If a field has no corresponding data in the text, omit it from the JSON or set it to null/empty as appropriate\n4. Follow the exact field names and structure from the schema below\n5. For arrays, only include items that have actual data from the text\n6. Use proper JSON data types (strings, numbers, booleans, arrays, objects)\n7. For dates, use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ) or null if not available\n8. For boolean fields, use true/false based on explicit information only",
					},
					{
						"type": "text",
						"text": "Resume to JSON Converter",
						"cache_control": map[string]interface{}{
							"type": "ephemeral",
						},
					},
				},
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_object",
			"schema": map[string]interface{}{
				"type":                 "object",
				"properties":           models.Resume{},
				"required":             []string{}, // Add required fields as needed
				"additionalProperties": false,
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("JSON marshal error: %v\n", err)
		return "Error preparing request"
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Request creation error: %v\n", err)
		return "Error creating request"
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Title", "Job-Applyer-Tool")
	req.Header.Set("HTTP-Referer", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request error: %v\n", err)
		return "Error making HTTP request"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP error: %d %s\n", resp.StatusCode, resp.Status)
		return "HTTP error from API"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response read error: %v\n", err)
		return "Error reading response"
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("JSON unmarshal error: %v\n", err)
		return "Error parsing response"
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Printf("No choices returned from ChatCompletion")
		return "No response from agent"
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"]

	var contentStr string
	switch v := content.(type) {
	case string:
		contentStr = v
	default:
		contentStr = fmt.Sprintf("%v", v)
	}

	fmt.Println(contentStr)
	return contentStr
}
