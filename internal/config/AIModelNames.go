package config

var modelNames = []string{
	"deepseek-r1-distill-llama-70b",
	"gemma2-9b-it",
	"llama-3.1-8b-instant",
	"llama-3.3-70b-versatile",
	"llama3-70b-8192",
	"mistral-saba-24b",
	"moonshotai/kimi-k2-instruct",
	"qwen/qwen3-32b",
}

func GetAllAIModelNames() []string {
	return modelNames
}
