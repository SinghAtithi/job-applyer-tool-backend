package database

import (
	_ "database/sql"
	"example.com/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	ConnectionString string
	Port             string
	Username         string
	Password         string
	DBName           string
	SSLMode          string
}

func getDBConnectionDetails() Database {
	// This function should retrieve the database connection details from a configuration file or environment variables.
	// For simplicity, we are returning a hardcoded example here.
	return Database{
		ConnectionString: "127.0.0.1",
		Port:             "5432",
		Username:         "ravan",
		Password:         "Abcd@1234",
		DBName:           "resumedb",
		SSLMode:          "disable",
	}
}

func NewDatabase() (*gorm.DB, error) {
	dataBaseConnectionDetails := getDBConnectionDetails()

	dsn := "host=" + dataBaseConnectionDetails.ConnectionString + " port=" + dataBaseConnectionDetails.Port + " user=" + dataBaseConnectionDetails.Username + " password=" + dataBaseConnectionDetails.Password + " dbname=" + dataBaseConnectionDetails.DBName + " sslmode=" + dataBaseConnectionDetails.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	// Perform any necessary setup or migrations here

	return db, nil
}

func EnsureTablesExist(db *gorm.DB) error {
	// AutoMigrate will create tables if they don't exist
	// and update schema if needed (adds new columns, indexes)
	err := db.AutoMigrate(
		&models.Resume{},
		&models.CoverLetterTable{},
		&models.JobDescriptionTable{},
	)

	if err != nil {
		return err
	}

	log.Println("All tables ensured to exist")
	return nil
}
