package publicPackage

import (
	"github.com/joho/godotenv"
)

func LoadDotEnvFile(fileName string) error {
	err := godotenv.Load(fileName)
	if err != nil {
		return err
	}
	return nil
}
