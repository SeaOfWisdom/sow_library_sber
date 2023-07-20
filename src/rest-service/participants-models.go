package rest

import (
	"fmt"

	"github.com/SeaOfWisdom/sow_library/src/service/storage"

	"github.com/ethereum/go-ethereum/common"
)

// NEW PARTICIPANT

type NewParticipantRequest struct {
	NickName    string `json:"nickname" example:"phd ***** destroyer"`
	Web3Address string `json:"web3_address"`
}

func (r *NewParticipantRequest) Validate() error {
	if r.NickName == "" {
		return fmt.Errorf("null nickname: %s", r.NickName)
	}

	if r.Web3Address == "" {
		return fmt.Errorf("null web3 address: %s", r.Web3Address)
	}

	if !common.IsHexAddress(r.Web3Address) {
		return fmt.Errorf("wrong web3 address: %s", r.Web3Address)
	}

	return nil
}

/// IfParticipantExists

type IfParticipantExistsResp struct {
	Result bool `json:"result" example:"true"`
}

// Updating a profile

type BasicInfoUpdateRequest struct {
	NickName string `json:"nickname" example:"phd ***** destroyer"`
}

func (r *BasicInfoUpdateRequest) Validate() error {
	return nil
}

/// ----- ----- Auth ----- -----

type AuthResp struct {
	Token    string                  `json:"jwt_token"`
	Role     storage.ParticipantRole `json:"role" example:"1"`
	NickName string                  `json:"nickname,omitempty" example:"phd ***** destroyer"`
}

// GET ROLE

type BasicInfo struct {
	NickName string                  `json:"nickname,omitempty" example:"phd ***** destroyer"`
	Role     storage.ParticipantRole `json:"role" example:"1"`
}
