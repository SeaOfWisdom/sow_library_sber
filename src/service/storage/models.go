package storage

import (
	"strings"
	"time"
)

// --- PostgreSQL ---

type ParticipantRole int

const (
	GuestRole     ParticipantRole = iota // 0 Guest
	ReaderRole                           // 1 Reader
	AuthorRole                           // 2 Author
	AdvisorRole                          // 3 Advisor
	ValidatorRole                        // 4 Validator
	AdminRole                            // 5 Admin
)

type WorkStatus string

var (
	PreReviewWorkStatus WorkStatus = "WORK_UNDER_PRE_REVIEW"
	ReviewWorkStatus    WorkStatus = "WORK_UNDER_REVIEW"
	OpenWorkStatus      WorkStatus = "WORK_OPEN"
	DeclinedWorkStatus  WorkStatus = "WORK_DECLINED"
)

type Participant struct {
	ID          string          `json:"-"`
	NickName    string          `gorm:"type:TEXT;uniqueIndex" json:"nickname"`
	Web3Address string          `gorm:"type:TEXT;uniqueIndex"`
	Role        ParticipantRole `json:"role,omitempty"`
	Language    string          `json:"language,omitempty"` // 'ru', 'en'
	CreatedAt   time.Time       `json:"-"`
}

type ParticipantsWork struct {
	ID            string     `json:"-"`
	ParticipantID string     `gorm:"type:TEXT" json:"-"`
	WorkID        string     `gorm:"type:TEXT" json:"-"`
	Tags          string     `gorm:"type:TEXT"`
	NFTAddress    string     `gorm:"type:TEXT" json:"nft_address"`
	Status        WorkStatus `json:"status,omitempty"`
	CreatedAt     time.Time  `json:"created_date,omitempty"`
}

func (w *ParticipantsWork) IsShow(participant *Participant, purchased bool) (work, content bool) {
	work = w.Status == OpenWorkStatus
	if participant != nil {
		work = work || w.ParticipantID == participant.ID || participant.Role >= ValidatorRole
		content = w.ParticipantID == participant.ID || participant.Role >= ValidatorRole
	}
	content = content || purchased

	return
}

type ParticipantsPurpose struct {
	ID            string    `json:"-"`
	ParticipantID string    `gorm:"type:TEXT" json:"-"`
	WorkID        string    `gorm:"type:TEXT" json:"-"`
	CreatedAt     time.Time `json:"created_date,omitempty"`
}

type ParticipantsBookmark struct {
	ID            string    `json:"-"`
	ParticipantID string    `gorm:"type:TEXT" json:"-"`
	WorkID        string    `gorm:"type:TEXT" json:"-"`
	CreatedAt     time.Time `json:"created_date,omitempty"`
}

type ParticipantsWorkReview struct {
	ID            string           `json:"-"`
	ParticipantID string           `gorm:"type:TEXT" json:"-"`
	WorkID        string           `gorm:"type:TEXT" json:"-"`
	Status        WorkReviewStatus `json:"status"`
	CreatedAt     time.Time        `gorm:"type:TIMESTAMP WITH TIME ZONE;default:now()" json:"created_date,omitempty,"`
	UpdatedAt     time.Time        `gorm:"type:TIMESTAMP WITH TIME ZONE;default:now()" json:"updated_date,omitempty"`
}

type AuthorResponse struct {
	BasicInfo  *Participant `json:"basic_info"`
	AuthorInfo *Author      `json:"author_info"`
}

type ValidatorResponse struct {
	BasicInfo     *Participant `json:"basic_info"`
	ValidatorInfo *Validator   `json:"validator_info"`
}

// --- Mongo ---

type Work struct {
	// BASE INFORMATION
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Annotation string     `json:"annotation"`
	AuthorID   string     `bson:"author_id" json:"-"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"-"`
	ReleasedAt time.Time  `bson:"released_at" json:"-"`
	Tags       []string   `jsob:"tags"`
	Price      string     `json:"price,omitempty"`
	Sources    string     `json:"sources,omitempty"`
	Language   string     `json:"language,omitempty"`
	Status     WorkStatus `bson:"status" json:"status,omitempty"`
	// BODY INFORMATION
	Content *WorkContent `json:"content"`
}

type WorkReviewStatus string

var (
	WorkReviewSkipped    WorkReviewStatus = "WORK_REVIEW_SKIPPED"
	WorkReviewInProgress WorkReviewStatus = "WORK_REVIEW_IN_PROGRESS"
	WorkReviewRejected   WorkReviewStatus = "WORK_REVIEW_DECLINED"
	WorkReviewSubmitted  WorkReviewStatus = "WORK_REVIEW_SUBMITTED"
	// after the admin's decision has been made.
	WorkReviewAccepted WorkReviewStatus = "WORK_REVIEW_ACCEPTED"
)

// StringToReviewStatus - convert string value to status. Default value = WORK_REVIEW_SUBMITTED
func StringToReviewStatus(val string) (out WorkReviewStatus) {
	out = WorkReviewSubmitted

	switch strings.ToUpper(val) {
	case "WORK_REVIEW_SKIPPED":
		out = WorkReviewSkipped
	case "WORK_REVIEW_ACCEPTED":
		out = WorkReviewAccepted
	case "WORK_REVIEW_DECLINED":
		out = WorkReviewRejected
	case "WORK_REVIEW_IN_PROGRESS":
		out = WorkReviewInProgress
	}

	return
}

type WorkReview struct {
	ID        string           `json:"id"`
	WorkID    string           `bson:"work_id" json:"work_id"`
	CreatedAt time.Time        `bson:"created_at" json:"created_date"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_date"`
	Language  string           `bson:"language" json:"language"`
	Status    WorkReviewStatus `bson:"status" json:"status"`
	// BODY REVIEW
	Body *WorkReviewBody `json:"body"`
}

type WorkReviewBody struct {
	Questionnaire *WorkReviewQuestionnaire `bson:"questionnaire" json:"questionnaire"`
	Review        string                   `json:"review"`
}
type WorkReviewQuestionnaire struct {
	Questions map[string]int64 `json:"questions"` // 0 - не согласен, 4 - согласен
}

type Validator struct {
	ID           string    `json:"-"` // postgresSQL id
	Name         string    `json:"name"`
	MiddleName   string    `json:"middlename,omitempty"`
	Surname      string    `json:"surname"`
	EmailAddress string    `bson:"email_address" json:"email_address"`
	Orcid        string    `json:"orcid,omitempty"`
	Sciences     []string  `json:"sciences,omitempty"`
	Language     string    `json:"language,omitempty"`
	DiplomaID    string    `bson:"diploma_id" json:"diploma_id,omitempty"` // referrenceKey
	CreatedAt    time.Time `bson:"created_at" json:"-"`
	UpdatedAt    time.Time `bson:"updated_at" json:"-"`
}

type DocumentType int

const (
	DiplomaType DocumentType = iota
)

type Document struct {
	ID      string
	Type    DocumentType `json:"type,omitempty"`
	Content interface{}  `json:"content"` // Diploma
}

type Diploma struct {
	ID            string `json:"-"` // postgresSQL id
	Degree        string
	Topics        []string
	DiplomaNumber int64
	DefenseNumber int64
	DefenseDate   time.Time
	OrderNumber   int64
	OrderDate     time.Time
	CreatedAt     time.Time `bson:"created_at" json:"-"`
	UpdatedAt     time.Time `bson:"updated_at" json:"-"`
}

type WorkContent struct {
	WorkData string `json:"work_data"`
}

type Author struct {
	ID                 string    `json:"-"` // postgresSQL id
	Name               string    `json:"name"`
	MiddleName         string    `json:"middlename,omitempty"`
	Surname            string    `json:"surname"`
	EmailAddress       string    `bson:"email_address" json:"email_address"`
	Orcid              string    `json:"orcid,omitempty"`
	Sciences           []string  `json:"sciences,omitempty"`
	Language           string    `json:"language,omitempty"`
	ScholarShipProfile string    `json:"scholar_ship_profile,omitempty"`
	CreatedAt          time.Time `bson:"created_at" json:"-"`
	UpdatedAt          time.Time `bson:"updated_at" json:"-"`
}

// WorkResponse consists of the information regarding the work and its author
type WorkResponse struct {
	Work       *Work           `json:"work"`
	Author     *AuthorResponse `json:"author_info"`
	Bookmarked bool            `json:"bookmarked"`
}
