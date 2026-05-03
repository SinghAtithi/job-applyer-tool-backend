package config

var modelNames = []string{
	"deepseek-r1-distill-llama-70b",
	"llama-3.3-70b-versatile",
	"meta-llama/llama-4-maverick-17b-128e-instruct",
	"openai/gpt-oss-120b",
	"qwen/qwen3-32b",
	"moonshotai/kimi-k2-instruct-0905",
}

func GetAllAIModelNames() []string {
	return modelNames
}
