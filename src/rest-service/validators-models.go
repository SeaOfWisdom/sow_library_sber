package rest

import (
	"fmt"

	"github.com/SeaOfWisdom/sow_library/src/service/storage"
)

type BecomeValidatorRequest struct {
	EmailAddress string   `json:"email_address"` // mandatory
	Name         string   `json:"name"`          // mandatory
	Surname      string   `json:"surname"`       // mandatory
	Middlename   string   `json:"middlename"`
	Orcid        string   `json:"orcid"`
	Sciences     []string `json:"sciences"`
	Language     string   `json:"language"` // mandatory
}

func (r *BecomeValidatorRequest) Validate() error {
	// TODO
	if r.EmailAddress == "" {
		return fmt.Errorf("wrong email address: %s", r.EmailAddress)
	}
	if r.Name == "" {
		return fmt.Errorf("wrong name: %s", r.Name)
	}
	if r.Surname == "" {
		return fmt.Errorf("wrong surname: %s", r.Surname)
	}
	if r.Language == "" {
		return fmt.Errorf("wrong language: %s", r.Language)
	}
	return nil
}

type UpdateValidatorRequest struct {
	EmailAddress string   `json:"email_address"`
	Name         string   `json:"name"`
	Surname      string   `json:"surname"`
	Middlename   string   `json:"middlename"`
	Orcid        string   `json:"orcid"`
	Sciences     []string `json:"sciences"`
	Language     string   `json:"language"`
}

func (r *UpdateValidatorRequest) Validate() error {
	return nil
}

type WorkReviewRequest struct {
	Review *storage.WorkReview `json:"review"`
}

func (r *WorkReviewRequest) Validate() error {
	if r.Review == nil {
		return fmt.Errorf("request is null: %v", r.Review)
	}
	if r.Review.WorkID == "" {
		return fmt.Errorf("work id is null: %s", r.Review.WorkID)
	}
	// if r.Review.ValidatorID == "" {
	// 	return fmt.Errorf("wrong validator id: %s", r.Review.ValidatorID)
	// }

	fmt.Println("REVIEW BODY: ", r.Review.Body.Questionnaire, r.Review.Body.Review)
	// if r.Review.Body == nil {
	// 	return fmt.Errorf("review body is null: %v", r.Review.Body)
	// }

	if r.Review.Body.Questionnaire == nil && r.Review.Body.Review == "" {
		return fmt.Errorf("review body is null: %v", r.Review.Body)
	}
	return nil
}
