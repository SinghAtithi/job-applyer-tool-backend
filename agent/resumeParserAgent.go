package agent

func GetResumeDetailInJsonFormat(resumeDetails string) string {
	// Convert the resume details to JSON format
	// This is a placeholder implementation; you would replace this with actual JSON conversion logic
	//jsonData := `{"resumeDetails": "` + resumeDetails + `"}`

	resumeDetails = MakeAgentCall(resumeDetails)

	// Return the JSON data as a string
	return resumeDetails
}
