package publicPackage

import (
	"github.com/joho/godotenv"
)

func LoadDotEnvFile(fileName string) error {
	return godotenv.Load(fileName)
}
