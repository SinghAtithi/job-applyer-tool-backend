// Package models provides data structures for the application's core domain models.
//
// This file defines the Resume struct and related types, representing a user's resume and its components.
// It includes embedded structs for personal information, education, experience, skills, projects, hobbies, and additional info.
//
// Author: Atithi (coder_ravan)
// Date: July 20, 2025
//
// Usage:
//   These models are used throughout the backend for resume management, parsing, and persistence.
//
// Struct Documentation:
//
// Resume: Represents a user's resume, including metadata and comprehensive resume data.
// PersonalInfo: Contains personal details of the user.
// Education: Represents an educational qualification (see below for struct definition).
// Experience: Represents a work experience entry (see below for struct definition).
// Skills: Represents a user's skills (see below for struct definition).
// Project: Represents a project entry (see below for struct definition).
// AdditionalInfo: Contains additional information about the user (see below for struct definition).
// Address: Represents a user's address (see below for struct definition).
//
// Functions in this file (if any) are documented individually below their definition.

package models

import (
	"errors"
	"fmt"
	"time"
)

// Resume represents a user's resume, including metadata and comprehensive resume data.
type Resume struct {
	ID          uint      `json:"id" gorm:"primaryKey"`                                                          // Unique identifier for the resume
	UserID      uint      `json:"user_id" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"` // Associated user ID
	Title       string    `json:"title"`                                                                         // Title of the resume
	Description string    `json:"description,omitempty"`                                                         // Optional description
	FileType    string    `json:"file_type,omitempty"`                                                           // File type (e.g., PDF, DOCX)
	IsActive    bool      `json:"is_active" gorm:"default:true"`                                                 // Indicates if the resume is active
	CreatedAt   time.Time `json:"created_at"`                                                                    // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at"`                                                                    // Last update timestamp

	// Embedded comprehensive resume data
	PersonalInfo   PersonalInfo    `json:"personal_info" gorm:"embedded;embeddedPrefix:personal_"`               // Personal information
	Summary        string          `json:"summary,omitempty" gorm:"embedded;embeddedPrefix:summary_"`            // Professional summary
	Education      []Education     `json:"education,omitempty" gorm:"serializer:json"`                           // List of education entries
	Experience     []Experience    `json:"experience,omitempty" gorm:"serializer:json"`                          // List of work experiences
	Skills         *Skills         `json:"skills,omitempty" gorm:"embedded;embeddedPrefix:skills_"`              // Skills
	Projects       []Project       `json:"projects,omitempty" gorm:"serializer:json"`                            // List of projects
	Hobbies        []string        `json:"hobbies,omitempty" gorm:"serializer:json"`                             // List of hobbies
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty" gorm:"embedded;embeddedPrefix:additional_"` // Additional information
}

// PersonalInfo contains personal details of the user.
type PersonalInfo struct {
	FirstName     string     `json:"first_name"`                                                // User's first name
	LastName      string     `json:"last_name"`                                                 // User's last name
	MiddleName    string     `json:"middle_name,omitempty"`                                     // Optional middle name
	Title         string     `json:"title,omitempty"`                                           // e.g., "Dr.", "Mr.", "Ms."
	Email         string     `json:"email"`                                                     // Email address
	Phone         string     `json:"phone,omitempty"`                                           // Optional phone number
	Address       *Address   `json:"address,omitempty" gorm:"embedded;embeddedPrefix:address_"` // Address
	Website       string     `json:"website,omitempty"`                                         // Optional website
	LinkedIn      string     `json:"linkedin,omitempty"`                                        // LinkedIn profile URL
	GitHub        string     `json:"github,omitempty"`                                          // GitHub profile URL
	Twitter       string     `json:"twitter,omitempty"`                                         // Twitter profile URL
	Portfolio     string     `json:"portfolio,omitempty"`                                       // Portfolio URL
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`                                   // Optional date of birth
	Nationality   string     `json:"nationality,omitempty"`                                     // Nationality
	Gender        string     `json:"gender,omitempty"`                                          // Gender
	MaritalStatus string     `json:"marital_status,omitempty"`                                  // Marital status
}

type Address struct {
	Street     string `json:"street,omitempty"`      // Street address
	City       string `json:"city,omitempty"`        // City
	State      string `json:"state,omitempty"`       // State or region
	Country    string `json:"country,omitempty"`     // Country
	PostalCode string `json:"postal_code,omitempty"` // Postal or ZIP code
}

// Education represents an educational qualification.
type Education struct {
	ID           uint       `json:"id,omitempty"`             // Unique identifier for the education entry
	Institution  string     `json:"institution"`              // Name of the institution
	Degree       string     `json:"degree"`                   // Degree or certification obtained
	FieldOfStudy string     `json:"field_of_study,omitempty"` // Field of study or major
	Grade        string     `json:"grade,omitempty"`          // GPA, CGPA, Percentage, etc.
	GradeType    string     `json:"grade_type,omitempty"`     // "GPA", "CGPA", "Percentage"
	StartDate    *time.Time `json:"start_date,omitempty"`     // Start date of the education
	EndDate      *time.Time `json:"end_date,omitempty"`       // End date of the education
	IsCurrent    bool       `json:"is_current,omitempty"`     // Indicates if it's the current education
}

// WorkExperience represents a work experience entry.
type Experience struct {
	ID               uint       `json:"id,omitempty"`                                       // Unique identifier for the experience entry
	JobTitle         string     `json:"job_title"`                                          // Job title
	Company          string     `json:"company"`                                            // Company name
	Role             string     `json:"role,omitempty"`                                     // Optional role description
	Location         string     `json:"location,omitempty"`                                 // Optional location
	StartDate        *time.Time `json:"start_date,omitempty"`                               // Start date of the experience
	EndDate          *time.Time `json:"end_date,omitempty"`                                 // End date of the experience
	IsCurrent        bool       `json:"is_current,omitempty"`                               // Indicates if it's the current job
	EmploymentType   string     `json:"employment_type,omitempty"`                          // Full-time, Part-time, Contract, Internship
	Description      string     `json:"description,omitempty"`                              // Job description
	Responsibilities []string   `json:"responsibilities,omitempty" gorm:"serializer:json"`  // List of responsibilities
	Achievements     []string   `json:"achievements,omitempty" gorm:"serializer:json"`      // List of achievements
	TechnologiesUsed []string   `json:"technologies_used,omitempty" gorm:"serializer:json"` // List of technologies used
}

// Skills represents a user's skills.
type Skills struct {
	Technical        []TechnicalSkill `json:"technical,omitempty" gorm:"serializer:json"`         // List of technical skills
	Tools            []string         `json:"tools,omitempty" gorm:"serializer:json"`             // List of tools
	Frameworks       []string         `json:"frameworks,omitempty" gorm:"serializer:json"`        // List of frameworks
	Libraries        []string         `json:"libraries,omitempty" gorm:"serializer:json"`         // List of libraries
	Databases        []string         `json:"databases,omitempty" gorm:"serializer:json"`         // List of databases
	CloudPlatforms   []string         `json:"cloud_platforms,omitempty" gorm:"serializer:json"`   // List of cloud platforms
	OperatingSystems []string         `json:"operating_systems,omitempty" gorm:"serializer:json"` // List of operating systems
	Methodologies    []string         `json:"methodologies,omitempty" gorm:"serializer:json"`     // List of methodologies
}

type TechnicalSkill struct {
	Name       string `json:"name"`                 // Name of the technical skill
	Level      string `json:"level,omitempty"`      // Beginner, Intermediate, Advanced, Expert
	Experience int    `json:"experience,omitempty"` // in years
	Category   string `json:"category,omitempty"`   // Programming, Database, Cloud, etc.
}

// Project represents a project entry.
type Project struct {
	ID           uint       `json:"id,omitempty"`                                  // Unique identifier for the project
	Name         string     `json:"name"`                                          // Project name
	Description  string     `json:"description,omitempty"`                         // Optional project description
	Technologies []string   `json:"technologies,omitempty" gorm:"serializer:json"` // List of technologies used
	URL          string     `json:"url,omitempty"`                                 // Optional project URL
	GitHubURL    string     `json:"github_url,omitempty"`                          // Optional GitHub URL
	StartDate    *time.Time `json:"start_date,omitempty"`                          // Start date of the project
	EndDate      *time.Time `json:"end_date,omitempty"`                            // End date of the project
	IsCurrent    bool       `json:"is_current,omitempty"`                          // Indicates if it's the current project
}

// AdditionalInfo contains additional information about the user.
type AdditionalInfo struct {
	Availability      string                 `json:"availability,omitempty"`                         // Availability status
	NoticePeriod      string                 `json:"notice_period,omitempty"`                        // Notice period duration
	Salary            string                 `json:"salary,omitempty"`                               // Expected salary
	WillingToRelocate bool                   `json:"willing_to_relocate,omitempty"`                  // Willingness to relocate
	WillingToTravel   string                 `json:"willing_to_travel,omitempty"`                    // Willingness to travel
	WorkAuthorization string                 `json:"work_authorization,omitempty"`                   // Work authorization status
	SecurityClearance string                 `json:"security_clearance,omitempty"`                   // Security clearance level
	DrivingLicense    string                 `json:"driving_license,omitempty"`                      // Driving license details
	MilitaryService   string                 `json:"military_service,omitempty"`                     // Military service details
	Disabilities      string                 `json:"disabilities,omitempty"`                         // Disability status
	EmergencyContact  string                 `json:"emergency_contact,omitempty"`                    // Emergency contact details
	PreferredLocation string                 `json:"preferred_location,omitempty"`                   // Preferred job location
	RemoteWork        bool                   `json:"remote_work,omitempty"`                          // Indicates if remote work is preferred
	CustomFields      map[string]interface{} `json:"custom_fields,omitempty" gorm:"serializer:json"` // Custom fields as key-value pairs
}

// Resume builder helper methods
func (r *Resume) GetFullName() string {
	if r.PersonalInfo.MiddleName != "" {
		return r.PersonalInfo.FirstName + " " + r.PersonalInfo.MiddleName + " " + r.PersonalInfo.LastName
	}
	return r.PersonalInfo.FirstName + " " + r.PersonalInfo.LastName
}

func (r *Resume) GetPrimaryContact() (string, string) {
	return r.PersonalInfo.Email, r.PersonalInfo.Phone
}

func (r *Resume) GetCurrentPosition() *Experience {
	for _, exp := range r.Experience {
		if exp.IsCurrent {
			return &exp
		}
	}
	return nil
}

func (r *Resume) GetLatestEducation() *Education {
	if len(r.Education) == 0 {
		return nil
	}

	latest := &r.Education[0]
	for _, edu := range r.Education {
		if edu.EndDate != nil && (latest.EndDate == nil || edu.EndDate.After(*latest.EndDate)) {
			latest = &edu
		}
	}
	return latest
}

func (r *Resume) GetTotalExperience() int {
	totalMonths := 0
	for _, exp := range r.Experience {
		if exp.StartDate != nil {
			endDate := time.Now()
			if exp.EndDate != nil {
				endDate = *exp.EndDate
			}
			months := int(endDate.Sub(*exp.StartDate).Hours() / (24 * 30))
			totalMonths += months
		}
	}
	return totalMonths / 12 // Convert to years
}

func (r *Resume) GetSkillsByCategory(category string) []string {
	if r.Skills == nil {
		return nil
	}

	var skills []string
	for _, skill := range r.Skills.Technical {
		if skill.Category == category {
			skills = append(skills, skill.Name)
		}
	}
	return skills
}

func (r *Resume) AddCustomField(key string, value interface{}) {
	if r.AdditionalInfo == nil {
		r.AdditionalInfo = &AdditionalInfo{}
	}
	if r.AdditionalInfo.CustomFields == nil {
		r.AdditionalInfo.CustomFields = make(map[string]interface{})
	}
	r.AdditionalInfo.CustomFields[key] = value
}

// Validation methods
func (r *Resume) Validate() error {
	if r.PersonalInfo.FirstName == "" {
		return errors.New("first name is required")
	}
	if r.PersonalInfo.LastName == "" {
		return errors.New("last name is required")
	}
	if r.PersonalInfo.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (e *Experience) GetDuration() string {
	if e.StartDate == nil {
		return ""
	}

	endDate := time.Now()
	if e.EndDate != nil {
		endDate = *e.EndDate
	}

	duration := endDate.Sub(*e.StartDate)
	years := int(duration.Hours() / (24 * 365))
	months := int(duration.Hours()/(24*30)) % 12

	if years > 0 && months > 0 {
		return fmt.Sprintf("%d years %d months", years, months)
	} else if years > 0 {
		return fmt.Sprintf("%d years", years)
	} else if months > 0 {
		return fmt.Sprintf("%d months", months)
	}
	return "Less than a month"
}

func (e *Education) GetDuration() string {
	if e.StartDate == nil || e.EndDate == nil {
		return ""
	}

	duration := e.EndDate.Sub(*e.StartDate)
	years := int(duration.Hours() / (24 * 365))

	if years > 0 {
		return fmt.Sprintf("%d years", years)
	}
	return "Less than a year"
}

func GetResumeSchemaSystemPrompt() string {
	return ResumeSchemaSystemPrompt
}

var ResumeSchemaSystemPrompt = "You are a resume data extraction system. Your task is to parse the provided resume text and return ONLY a JSON object that matches the following schema. Follow these rules strictly:\n\n1. Return ONLY valid JSON - no additional text, explanations, or formatting\n2. Use only the data present in the provided text - do not create, assume, or invent any information\n3. If a field has no corresponding data in the text, omit it from the JSON or set it to null/empty as appropriate\n4. Follow the exact field names and structure from the schema below\n5. For arrays, only include items that have actual data from the text\n6. Use proper JSON data types (strings, numbers, booleans, arrays, objects)\n7. For dates, use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ) or null if not available\n8. For boolean fields, use true/false based on explicit information only\n\nJSON Schema to follow:\n\n```json\n{\n  \"id\": number,\n  \"user_id\": number,\n  \"title\": string,\n  \"description\": string,\n  \"file_type\": string,\n  \"is_active\": boolean,\n  \"created_at\": string, // ISO 8601 datetime format\n  \"updated_at\": string, // ISO 8601 datetime format\n  \n  \"personal_info\": {\n    \"first_name\": string,\n    \"last_name\": string,\n    \"middle_name\": string,\n    \"title\": string, // e.g., \"Dr.\", \"Mr.\", \"Ms.\"\n    \"email\": string,\n    \"phone\": string,\n    \"alternate_phone\": string,\n    \"address\": {\n      \"street\": string,\n      \"city\": string,\n      \"state\": string,\n      \"country\": string,\n      \"postal_code\": string,\n      \"is_present\": boolean\n    },\n    \"website\": string,\n    \"linkedin\": string,\n    \"github\": string,\n    \"twitter\": string,\n    \"instagram\": string,\n    \"facebook\": string,\n    \"portfolio\": string,\n    \"date_of_birth\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n    \"nationality\": string,\n    \"gender\": string,\n    \"marital_status\": string,\n    \"profile_picture\": string,\n    \"social_links\": object // map of string to string\n  },\n  \n  \"summary\": {\n    \"content\": string,\n    \"keywords\": array, // array of strings\n    \"years_of_experience\": number\n  },\n  \n  \"objective\": {\n    \"content\": string,\n    \"career_goal\": string,\n    \"industry\": string\n  },\n  \n  \"education\": [ // array of objects\n    {\n      \"id\": number,\n      \"institution\": string,\n      \"degree\": string,\n      \"field_of_study\": string,\n      \"grade\": string, // GPA, CGPA, Percentage, etc.\n      \"grade_type\": string, // \"GPA\", \"CGPA\", \"Percentage\"\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"location\": string,\n      \"description\": string,\n      \"achievements\": array, // array of strings\n      \"coursework\": array, // array of strings\n      \"honors\": array, // array of strings\n      \"thesis\": string,\n      \"advisor\": string,\n      \"activities\": array // array of strings\n    }\n  ],\n  \n  \"experience\": [ // array of objects\n    {\n      \"id\": number,\n      \"job_title\": string,\n      \"company\": string,\n      \"department\": string,\n      \"location\": string,\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"employment_type\": string, // Full-time, Part-time, Contract, Internship\n      \"description\": string,\n      \"responsibilities\": array, // array of strings\n      \"achievements\": array, // array of strings\n      \"technologies_used\": array, // array of strings\n      \"supervisor\": string,\n      \"salary\": string,\n      \"reason_for_leaving\": string,\n      \"skills\": array, // array of strings\n      \"projects\": array // array of strings\n    }\n  ],\n  \n  \"skills\": {\n    \"technical\": [ // array of objects\n      {\n        \"name\": string,\n        \"level\": string, // Beginner, Intermediate, Advanced, Expert\n        \"years_of_experience\": number,\n        \"category\": string // Programming, Database, Cloud, etc.\n      }\n    ],\n    \"soft\": [ // array of objects\n      {\n        \"name\": string,\n        \"level\": string,\n        \"description\": string\n      }\n    ],\n    \"languages\": [ // array of objects\n      {\n        \"name\": string,\n        \"proficiency\": string, // Native, Fluent, Conversational, Basic\n        \"speaking\": string,\n        \"writing\": string,\n        \"reading\": string\n      }\n    ],\n    \"tools\": array, // array of strings\n    \"frameworks\": array, // array of strings\n    \"libraries\": array, // array of strings\n    \"databases\": array, // array of strings\n    \"cloud_platforms\": array, // array of strings\n    \"operating_systems\": array, // array of strings\n    \"methodologies\": array, // array of strings\n    \"industry\": array // array of strings\n  },\n  \n  \"projects\": [ // array of objects\n    {\n      \"id\": number,\n      \"name\": string,\n      \"description\": string,\n      \"role\": string,\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"technologies\": array, // array of strings\n      \"features\": array, // array of strings\n      \"achievements\": array, // array of strings\n      \"url\": string,\n      \"github_url\": string,\n      \"demo_url\": string,\n      \"category\": string, // Web, Mobile, Desktop, AI/ML, etc.\n      \"status\": string, // Completed, In Progress, Planned\n      \"team_size\": number,\n      \"client\": string,\n      \"budget\": string,\n      \"challenges\": array, // array of strings\n      \"learnings\": array // array of strings\n    }\n  ],\n  \n  \"certifications\": [ // array of objects\n    {\n      \"id\": number,\n      \"name\": string,\n      \"issuing_organization\": string,\n      \"issue_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"expiry_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"credential_id\": string,\n      \"url\": string,\n      \"description\": string,\n      \"skills\": array, // array of strings\n      \"verification_url\": string,\n      \"score\": string,\n      \"is_active\": boolean\n    }\n  ],\n  \n  \"publications\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"authors\": array, // array of strings\n      \"journal\": string,\n      \"publisher\": string,\n      \"publication_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"volume\": string,\n      \"issue\": string,\n      \"pages\": string,\n      \"doi\": string,\n      \"url\": string,\n      \"abstract\": string,\n      \"keywords\": array, // array of strings\n      \"citation_count\": number,\n      \"type\": string, // Journal, Conference, Book, etc.\n      \"status\": string // Published, Under Review, Draft\n    }\n  ],\n  \n  \"awards\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"issuing_organization\": string,\n      \"date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"description\": string,\n      \"category\": string, // Academic, Professional, Competition\n      \"level\": string, // International, National, Regional, Local\n      \"rank\": string,\n      \"prize\": string,\n      \"url\": string,\n      \"skills\": array // array of strings\n    }\n  ],\n  \n  \"languages\": [ // array of objects\n    {\n      \"id\": number,\n      \"name\": string,\n      \"proficiency\": string, // Native, Fluent, Conversational, Basic\n      \"speaking\": string,\n      \"writing\": string,\n      \"reading\": string,\n      \"listening\": string,\n      \"certification\": string,\n      \"test_score\": string,\n      \"years_of_experience\": number\n    }\n  ],\n  \n  \"volunteer\": [ // array of objects\n    {\n      \"id\": number,\n      \"organization\": string,\n      \"role\": string,\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"location\": string,\n      \"description\": string,\n      \"achievements\": array, // array of strings\n      \"skills\": array, // array of strings\n      \"hours_per_week\": number,\n      \"total_hours\": number,\n      \"cause\": string,\n      \"website\": string\n    }\n  ],\n  \n  \"leadership\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"organization\": string,\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"location\": string,\n      \"description\": string,\n      \"responsibilities\": array, // array of strings\n      \"achievements\": array, // array of strings\n      \"team_size\": number,\n      \"budget\": string,\n      \"skills\": array, // array of strings\n      \"type\": string // Professional, Academic, Community\n    }\n  ],\n  \n  \"references\": [ // array of objects\n    {\n      \"id\": number,\n      \"name\": string,\n      \"title\": string,\n      \"company\": string,\n      \"email\": string,\n      \"phone\": string,\n      \"relationship\": string,\n      \"years_known\": number,\n      \"linkedin\": string,\n      \"can_contact\": boolean,\n      \"note\": string\n    }\n  ],\n  \n  \"courses\": [ // array of objects\n    {\n      \"id\": number,\n      \"name\": string,\n      \"provider\": string,\n      \"instructor\": string,\n      \"completion_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"duration\": string,\n      \"grade\": string,\n      \"credential_id\": string,\n      \"url\": string,\n      \"description\": string,\n      \"skills\": array, // array of strings\n      \"category\": string,\n      \"is_online\": boolean,\n      \"certificate\": string\n    }\n  ],\n  \n  \"research\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"institution\": string,\n      \"supervisor\": string,\n      \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"end_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"is_current\": boolean,\n      \"field\": string,\n      \"description\": string,\n      \"methodology\": array, // array of strings\n      \"findings\": array, // array of strings\n      \"publications\": array, // array of strings\n      \"technologies\": array, // array of strings\n      \"funding\": string,\n      \"collaborators_count\": number\n    }\n  ],\n  \n  \"patents\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"patent_number\": string,\n      \"application_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"issue_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"status\": string, // Pending, Granted, Expired\n      \"inventors\": array, // array of strings\n      \"assignee\": string,\n      \"country\": string,\n      \"description\": string,\n      \"claims\": array, // array of strings\n      \"url\": string,\n      \"field\": string\n    }\n  ],\n  \n  \"conferences\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"event\": string,\n      \"location\": string,\n      \"date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"type\": string, // Presenter, Attendee, Organizer\n      \"presentation_type\": string, // Oral, Poster, Workshop\n      \"abstract\": string,\n      \"audience\": string,\n      \"url\": string,\n      \"co_authors\": array, // array of strings\n      \"awards\": array // array of strings\n    }\n  ],\n  \n  \"portfolio\": [ // array of objects\n    {\n      \"id\": number,\n      \"title\": string,\n      \"type\": string, // Image, Video, Document, Website, etc.\n      \"url\": string,\n      \"thumbnail_url\": string,\n      \"description\": string,\n      \"category\": string,\n      \"tags\": array, // array of strings\n      \"date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n      \"client\": string,\n      \"technologies\": array, // array of strings\n      \"is_public\": boolean\n    }\n  ],\n  \n  \"hobbies\": array, // array of strings\n  \n  \"additional_info\": {\n    \"availability\": string,\n    \"notice_period\": string,\n    \"salary\": string,\n    \"willing_to_relocate\": boolean,\n    \"willing_to_travel\": string,\n    \"work_authorization\": string,\n    \"security_clearance\": string,\n    \"driving_license\": string,\n    \"military_service\": string,\n    \"disabilities\": string,\n    \"emergency_contact\": string,\n    \"preferred_location\": string,\n    \"remote_work\": boolean,\n    \"part_time\": boolean,\n    \"freelance\": boolean,\n    \"start_date\": string, // ISO 8601 datetime format or null (2006-01-02T15:04:05Z07:00 should look like this format)\n    \"personal_statement\": string,\n    \"career_goals\": string,\n    \"professional_memberships\": array, // array of strings\n    \"board_positions\": array, // array of strings\n    \"media_appearances\": array, // array of strings\n    \"speaking_engagements\": array, // array of strings\n    \"test_scores\": object, // map of string to string\n    \"custom_fields\": object // map of string to any type\n  }\n}\n```\n\nExtract and structure the data exactly as it appears in the text, maintaining accuracy and completeness while following this JSON schema structure. Return ONLY the JSON object in stringify format without any code block."
