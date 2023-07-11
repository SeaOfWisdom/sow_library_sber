package rest

import "fmt"

// BecomeAuthorRequest model info
// @Description User account information
// @Description with user id and username

type BecomeAuthorRequest struct {
	EmailAddress       string   `json:"email_address" example:"mr_math_phd@science.com"`
	Name               string   `json:"name"`
	Surname            string   `json:"surname"`
	Middlename         string   `json:"middlename"`
	Orcid              string   `json:"orcid"`
	Sciences           []string `json:"sciences"`
	Language           string   `json:"language"`
	ScholarShipProfile string   `json:"scholar_ship_profile"`
}

func (r *BecomeAuthorRequest) Validate() error {
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
	return nil
}

type BecomeAuthorDataResp struct {
	Sciences []string `json:"sciences"`
}

type UpdateAuthorRequest struct {
	EmailAddress       string   `json:"email_address"`
	Name               string   `json:"name"`
	Surname            string   `json:"surname"`
	Middlename         string   `json:"middlename"`
	Orcid              string   `json:"orcid"`
	Sciences           []string `json:"sciences"`
	Language           string   `json:"language"`
	ScholarShipProfile string   `json:"scholar_ship_profile"`
}

func (r *UpdateAuthorRequest) Validate() error {
	return nil
}
