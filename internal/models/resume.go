package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Base Resume struct
type Resume struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	FileType    string    `json:"file_type,omitempty"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Embedded comprehensive resume data
	PersonalInfo   PersonalInfo    `json:"personal_info" gorm:"embedded;embeddedPrefix:personal_"`
	Summary        *Summary        `json:"summary,omitempty" gorm:"embedded;embeddedPrefix:summary_"`
	Objective      *Objective      `json:"objective,omitempty" gorm:"embedded;embeddedPrefix:objective_"`
	Education      []Education     `json:"education,omitempty" gorm:"serializer:json"`
	Experience     []Experience    `json:"experience,omitempty" gorm:"serializer:json"`
	Skills         *Skills         `json:"skills,omitempty" gorm:"embedded;embeddedPrefix:skills_"`
	Projects       []Project       `json:"projects,omitempty" gorm:"serializer:json"`
	Certifications []Certification `json:"certifications,omitempty" gorm:"serializer:json"`
	Publications   []Publication   `json:"publications,omitempty" gorm:"serializer:json"`
	Awards         []Award         `json:"awards,omitempty" gorm:"serializer:json"`
	Languages      []Language      `json:"languages,omitempty" gorm:"serializer:json"`
	Volunteer      []VolunteerWork `json:"volunteer,omitempty" gorm:"serializer:json"`
	Leadership     []Leadership    `json:"leadership,omitempty" gorm:"serializer:json"`
	References     []Reference     `json:"references,omitempty" gorm:"serializer:json"`
	Courses        []Course        `json:"courses,omitempty" gorm:"serializer:json"`
	Research       []Research      `json:"research,omitempty" gorm:"serializer:json"`
	Patents        []Patent        `json:"patents,omitempty" gorm:"serializer:json"`
	Conferences    []Conference    `json:"conferences,omitempty" gorm:"serializer:json"`
	Portfolio      []PortfolioItem `json:"portfolio,omitempty" gorm:"serializer:json"`
	Hobbies        []string        `json:"hobbies,omitempty" gorm:"serializer:json"`
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty" gorm:"embedded;embeddedPrefix:additional_"`
}

// Personal Information
type PersonalInfo struct {
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	MiddleName     string      `json:"middle_name,omitempty"`
	Title          string      `json:"title,omitempty"` // e.g., "Dr.", "Mr.", "Ms."
	Email          string      `json:"email"`
	Phone          string      `json:"phone,omitempty"`
	AlternatePhone string      `json:"alternate_phone,omitempty"`
	Address        *Address    `json:"address,omitempty" gorm:"embedded;embeddedPrefix:address_"`
	Website        string      `json:"website,omitempty"`
	LinkedIn       string      `json:"linkedin,omitempty"`
	GitHub         string      `json:"github,omitempty"`
	Twitter        string      `json:"twitter,omitempty"`
	Instagram      string      `json:"instagram,omitempty"`
	Facebook       string      `json:"facebook,omitempty"`
	Portfolio      string      `json:"portfolio,omitempty"`
	DateOfBirth    *time.Time  `json:"date_of_birth,omitempty"`
	Nationality    string      `json:"nationality,omitempty"`
	Gender         string      `json:"gender,omitempty"`
	MaritalStatus  string      `json:"marital_status,omitempty"`
	ProfilePicture string      `json:"profile_picture,omitempty"`
	SocialLinks    SocialLinks `json:"social_links,omitempty" gorm:"serializer:json"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	IsPresent  bool   `json:"is_present,omitempty"`
}

type SocialLinks map[string]string // For flexible social media links

// Summary/Profile Section
type Summary struct {
	Content    string   `json:"content"`
	Keywords   []string `json:"keywords,omitempty" gorm:"serializer:json"`
	YearsOfExp int      `json:"years_of_experience,omitempty"`
}

// Objective Section
type Objective struct {
	Content    string `json:"content"`
	CareerGoal string `json:"career_goal,omitempty"`
	Industry   string `json:"industry,omitempty"`
}

// Education
type Education struct {
	ID           uint       `json:"id,omitempty"`
	Institution  string     `json:"institution"`
	Degree       string     `json:"degree"`
	FieldOfStudy string     `json:"field_of_study,omitempty"`
	Grade        string     `json:"grade,omitempty"`      // GPA, CGPA, Percentage, etc.
	GradeType    string     `json:"grade_type,omitempty"` // "GPA", "CGPA", "Percentage"
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	IsCurrent    bool       `json:"is_current,omitempty"`
	Location     string     `json:"location,omitempty"`
	Description  string     `json:"description,omitempty"`
	Achievements []string   `json:"achievements,omitempty" gorm:"serializer:json"`
	Coursework   []string   `json:"coursework,omitempty" gorm:"serializer:json"`
	Honors       []string   `json:"honors,omitempty" gorm:"serializer:json"`
	Thesis       string     `json:"thesis,omitempty"`
	Advisor      string     `json:"advisor,omitempty"`
	Activities   []string   `json:"activities,omitempty" gorm:"serializer:json"`
}

// Work Experience
type Experience struct {
	ID               uint       `json:"id,omitempty"`
	JobTitle         string     `json:"job_title"`
	Company          string     `json:"company"`
	Department       string     `json:"department,omitempty"`
	Location         string     `json:"location,omitempty"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	IsCurrent        bool       `json:"is_current,omitempty"`
	EmploymentType   string     `json:"employment_type,omitempty"` // Full-time, Part-time, Contract, Internship
	Description      string     `json:"description,omitempty"`
	Responsibilities []string   `json:"responsibilities,omitempty" gorm:"serializer:json"`
	Achievements     []string   `json:"achievements,omitempty" gorm:"serializer:json"`
	TechnologiesUsed []string   `json:"technologies_used,omitempty" gorm:"serializer:json"`
	Supervisor       string     `json:"supervisor,omitempty"`
	Salary           string     `json:"salary,omitempty"`
	ReasonForLeaving string     `json:"reason_for_leaving,omitempty"`
	Skills           []string   `json:"skills,omitempty" gorm:"serializer:json"`
	Projects         []string   `json:"projects,omitempty" gorm:"serializer:json"`
}

// Skills Section
type Skills struct {
	Technical        []TechnicalSkill `json:"technical,omitempty" gorm:"serializer:json"`
	Soft             []SoftSkill      `json:"soft,omitempty" gorm:"serializer:json"`
	Languages        []LanguageSkill  `json:"languages,omitempty" gorm:"serializer:json"`
	Tools            []string         `json:"tools,omitempty" gorm:"serializer:json"`
	Frameworks       []string         `json:"frameworks,omitempty" gorm:"serializer:json"`
	Libraries        []string         `json:"libraries,omitempty" gorm:"serializer:json"`
	Databases        []string         `json:"databases,omitempty" gorm:"serializer:json"`
	CloudPlatforms   []string         `json:"cloud_platforms,omitempty" gorm:"serializer:json"`
	OperatingSystems []string         `json:"operating_systems,omitempty" gorm:"serializer:json"`
	Methodologies    []string         `json:"methodologies,omitempty" gorm:"serializer:json"`
	Industry         []string         `json:"industry,omitempty" gorm:"serializer:json"`
}

type TechnicalSkill struct {
	Name       string `json:"name"`
	Level      string `json:"level,omitempty"` // Beginner, Intermediate, Advanced, Expert
	YearsOfExp int    `json:"years_of_experience,omitempty"`
	Category   string `json:"category,omitempty"` // Programming, Database, Cloud, etc.
}

type SoftSkill struct {
	Name        string `json:"name"`
	Level       string `json:"level,omitempty"`
	Description string `json:"description,omitempty"`
}

type LanguageSkill struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency,omitempty"` // Native, Fluent, Conversational, Basic
	Speaking    string `json:"speaking,omitempty"`
	Writing     string `json:"writing,omitempty"`
	Reading     string `json:"reading,omitempty"`
}

// Projects
type Project struct {
	ID           uint       `json:"id,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Role         string     `json:"role,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	IsCurrent    bool       `json:"is_current,omitempty"`
	Technologies []string   `json:"technologies,omitempty" gorm:"serializer:json"`
	Features     []string   `json:"features,omitempty" gorm:"serializer:json"`
	Achievements []string   `json:"achievements,omitempty" gorm:"serializer:json"`
	URL          string     `json:"url,omitempty"`
	GitHubURL    string     `json:"github_url,omitempty"`
	DemoURL      string     `json:"demo_url,omitempty"`
	Category     string     `json:"category,omitempty"` // Web, Mobile, Desktop, AI/ML, etc.
	Status       string     `json:"status,omitempty"`   // Completed, In Progress, Planned
	TeamSize     int        `json:"team_size,omitempty"`
	Client       string     `json:"client,omitempty"`
	Budget       string     `json:"budget,omitempty"`
	Challenges   []string   `json:"challenges,omitempty" gorm:"serializer:json"`
	Learnings    []string   `json:"learnings,omitempty" gorm:"serializer:json"`
}

// Certifications
type Certification struct {
	ID              uint       `json:"id,omitempty"`
	Name            string     `json:"name"`
	IssuingOrg      string     `json:"issuing_organization"`
	IssueDate       *time.Time `json:"issue_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CredentialID    string     `json:"credential_id,omitempty"`
	URL             string     `json:"url,omitempty"`
	Description     string     `json:"description,omitempty"`
	Skills          []string   `json:"skills,omitempty" gorm:"serializer:json"`
	VerificationURL string     `json:"verification_url,omitempty"`
	Score           string     `json:"score,omitempty"`
	IsActive        bool       `json:"is_active,omitempty"`
}

// Publications
type Publication struct {
	ID              uint       `json:"id,omitempty"`
	Title           string     `json:"title"`
	Authors         []string   `json:"authors,omitempty" gorm:"serializer:json"`
	Journal         string     `json:"journal,omitempty"`
	Publisher       string     `json:"publisher,omitempty"`
	PublicationDate *time.Time `json:"publication_date,omitempty"`
	Volume          string     `json:"volume,omitempty"`
	Issue           string     `json:"issue,omitempty"`
	Pages           string     `json:"pages,omitempty"`
	DOI             string     `json:"doi,omitempty"`
	URL             string     `json:"url,omitempty"`
	Abstract        string     `json:"abstract,omitempty"`
	Keywords        []string   `json:"keywords,omitempty" gorm:"serializer:json"`
	CitationCount   int        `json:"citation_count,omitempty"`
	Type            string     `json:"type,omitempty"`   // Journal, Conference, Book, etc.
	Status          string     `json:"status,omitempty"` // Published, Under Review, Draft
}

// Awards and Honors
type Award struct {
	ID          uint       `json:"id,omitempty"`
	Title       string     `json:"title"`
	IssuingOrg  string     `json:"issuing_organization"`
	Date        *time.Time `json:"date,omitempty"`
	Description string     `json:"description,omitempty"`
	Category    string     `json:"category,omitempty"` // Academic, Professional, Competition
	Level       string     `json:"level,omitempty"`    // International, National, Regional, Local
	Rank        string     `json:"rank,omitempty"`
	Prize       string     `json:"prize,omitempty"`
	URL         string     `json:"url,omitempty"`
	Skills      []string   `json:"skills,omitempty" gorm:"serializer:json"`
}

// Languages
type Language struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name"`
	Proficiency   string `json:"proficiency"` // Native, Fluent, Conversational, Basic
	Speaking      string `json:"speaking,omitempty"`
	Writing       string `json:"writing,omitempty"`
	Reading       string `json:"reading,omitempty"`
	Listening     string `json:"listening,omitempty"`
	Certification string `json:"certification,omitempty"`
	TestScore     string `json:"test_score,omitempty"`
	YearsOfExp    int    `json:"years_of_experience,omitempty"`
}

// Volunteer Work
type VolunteerWork struct {
	ID           uint       `json:"id,omitempty"`
	Organization string     `json:"organization"`
	Role         string     `json:"role"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	IsCurrent    bool       `json:"is_current,omitempty"`
	Location     string     `json:"location,omitempty"`
	Description  string     `json:"description,omitempty"`
	Achievements []string   `json:"achievements,omitempty" gorm:"serializer:json"`
	Skills       []string   `json:"skills,omitempty" gorm:"serializer:json"`
	HoursPerWeek int        `json:"hours_per_week,omitempty"`
	TotalHours   int        `json:"total_hours,omitempty"`
	Cause        string     `json:"cause,omitempty"`
	Website      string     `json:"website,omitempty"`
}

// Leadership Experience
type Leadership struct {
	ID               uint       `json:"id,omitempty"`
	Title            string     `json:"title"`
	Organization     string     `json:"organization"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	IsCurrent        bool       `json:"is_current,omitempty"`
	Location         string     `json:"location,omitempty"`
	Description      string     `json:"description,omitempty"`
	Responsibilities []string   `json:"responsibilities,omitempty" gorm:"serializer:json"`
	Achievements     []string   `json:"achievements,omitempty" gorm:"serializer:json"`
	TeamSize         int        `json:"team_size,omitempty"`
	Budget           string     `json:"budget,omitempty"`
	Skills           []string   `json:"skills,omitempty" gorm:"serializer:json"`
	Type             string     `json:"type,omitempty"` // Professional, Academic, Community
}

// References
type Reference struct {
	ID           uint   `json:"id,omitempty"`
	Name         string `json:"name"`
	Title        string `json:"title,omitempty"`
	Company      string `json:"company,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	YearsKnown   int    `json:"years_known,omitempty"`
	LinkedIn     string `json:"linkedin,omitempty"`
	CanContact   bool   `json:"can_contact,omitempty"`
	Note         string `json:"note,omitempty"`
}

// Courses and Training
type Course struct {
	ID             uint       `json:"id,omitempty"`
	Name           string     `json:"name"`
	Provider       string     `json:"provider,omitempty"`
	Instructor     string     `json:"instructor,omitempty"`
	CompletionDate *time.Time `json:"completion_date,omitempty"`
	Duration       string     `json:"duration,omitempty"`
	Grade          string     `json:"grade,omitempty"`
	CredentialID   string     `json:"credential_id,omitempty"`
	URL            string     `json:"url,omitempty"`
	Description    string     `json:"description,omitempty"`
	Skills         []string   `json:"skills,omitempty" gorm:"serializer:json"`
	Category       string     `json:"category,omitempty"`
	IsOnline       bool       `json:"is_online,omitempty"`
	Certificate    string     `json:"certificate,omitempty"`
}

// Research Experience
type Research struct {
	ID                 uint       `json:"id,omitempty"`
	Title              string     `json:"title"`
	Institution        string     `json:"institution"`
	Supervisor         string     `json:"supervisor,omitempty"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsCurrent          bool       `json:"is_current,omitempty"`
	Field              string     `json:"field,omitempty"`
	Description        string     `json:"description,omitempty"`
	Methodology        []string   `json:"methodology,omitempty" gorm:"serializer:json"`
	Findings           []string   `json:"findings,omitempty" gorm:"serializer:json"`
	Publications       []string   `json:"publications,omitempty" gorm:"serializer:json"`
	Technologies       []string   `json:"technologies,omitempty" gorm:"serializer:json"`
	Funding            string     `json:"funding,omitempty"`
	CollaboratorsCount int        `json:"collaborators_count,omitempty"`
}

// Patents
type Patent struct {
	ID              uint       `json:"id,omitempty"`
	Title           string     `json:"title"`
	PatentNumber    string     `json:"patent_number,omitempty"`
	ApplicationDate *time.Time `json:"application_date,omitempty"`
	IssueDate       *time.Time `json:"issue_date,omitempty"`
	Status          string     `json:"status,omitempty"` // Pending, Granted, Expired
	Inventors       []string   `json:"inventors,omitempty" gorm:"serializer:json"`
	Assignee        string     `json:"assignee,omitempty"`
	Country         string     `json:"country,omitempty"`
	Description     string     `json:"description,omitempty"`
	Claims          []string   `json:"claims,omitempty" gorm:"serializer:json"`
	URL             string     `json:"url,omitempty"`
	Field           string     `json:"field,omitempty"`
}

// Conferences and Presentations
type Conference struct {
	ID               uint       `json:"id,omitempty"`
	Title            string     `json:"title"`
	Event            string     `json:"event"`
	Location         string     `json:"location,omitempty"`
	Date             *time.Time `json:"date,omitempty"`
	Type             string     `json:"type,omitempty"`              // Presenter, Attendee, Organizer
	PresentationType string     `json:"presentation_type,omitempty"` // Oral, Poster, Workshop
	Abstract         string     `json:"abstract,omitempty"`
	Audience         string     `json:"audience,omitempty"`
	URL              string     `json:"url,omitempty"`
	CoAuthors        []string   `json:"co_authors,omitempty" gorm:"serializer:json"`
	Awards           []string   `json:"awards,omitempty" gorm:"serializer:json"`
}

// Portfolio Items
type PortfolioItem struct {
	ID           uint       `json:"id,omitempty"`
	Title        string     `json:"title"`
	Type         string     `json:"type,omitempty"` // Image, Video, Document, Website, etc.
	URL          string     `json:"url,omitempty"`
	ThumbnailURL string     `json:"thumbnail_url,omitempty"`
	Description  string     `json:"description,omitempty"`
	Category     string     `json:"category,omitempty"`
	Tags         []string   `json:"tags,omitempty" gorm:"serializer:json"`
	Date         *time.Time `json:"date,omitempty"`
	Client       string     `json:"client,omitempty"`
	Technologies []string   `json:"technologies,omitempty" gorm:"serializer:json"`
	IsPublic     bool       `json:"is_public,omitempty"`
}

// Additional Information
type AdditionalInfo struct {
	Availability            string                 `json:"availability,omitempty"`
	NoticePeriod            string                 `json:"notice_period,omitempty"`
	Salary                  string                 `json:"salary,omitempty"`
	WillingToRelocate       bool                   `json:"willing_to_relocate,omitempty"`
	WillingToTravel         string                 `json:"willing_to_travel,omitempty"`
	WorkAuthorization       string                 `json:"work_authorization,omitempty"`
	SecurityClearance       string                 `json:"security_clearance,omitempty"`
	DrivingLicense          string                 `json:"driving_license,omitempty"`
	MilitaryService         string                 `json:"military_service,omitempty"`
	Disabilities            string                 `json:"disabilities,omitempty"`
	EmergencyContact        string                 `json:"emergency_contact,omitempty"`
	PreferredLocation       string                 `json:"preferred_location,omitempty"`
	RemoteWork              bool                   `json:"remote_work,omitempty"`
	PartTime                bool                   `json:"part_time,omitempty"`
	Freelance               bool                   `json:"freelance,omitempty"`
	StartDate               *time.Time             `json:"start_date,omitempty"`
	PersonalStatement       string                 `json:"personal_statement,omitempty"`
	CareerGoals             string                 `json:"career_goals,omitempty"`
	ProfessionalMemberships []string               `json:"professional_memberships,omitempty" gorm:"serializer:json"`
	BoardPositions          []string               `json:"board_positions,omitempty" gorm:"serializer:json"`
	MediaAppearances        []string               `json:"media_appearances,omitempty" gorm:"serializer:json"`
	SpeakingEngagements     []string               `json:"speaking_engagements,omitempty" gorm:"serializer:json"`
	TestScores              map[string]string      `json:"test_scores,omitempty" gorm:"serializer:json"`
	CustomFields            map[string]interface{} `json:"custom_fields,omitempty" gorm:"serializer:json"`
}

// Helper methods for JSON serialization with GORM
func (s SocialLinks) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *SocialLinks) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-bytes value into SocialLinks")
	}

	return json.Unmarshal(bytes, s)
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
