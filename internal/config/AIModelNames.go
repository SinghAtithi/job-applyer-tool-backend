package config

var modelNames = []string{
	"deepseek-r1-distill-llama-70b",
	"llama-3.3-70b-versatile",
	"llama3-70b-8192",
	"moonshotai/kimi-k2-instruct",
	"qwen/qwen3-32b",
}

func GetAllAIModelNames() []string {
	return modelNames
}
