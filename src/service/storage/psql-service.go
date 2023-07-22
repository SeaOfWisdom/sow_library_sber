package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Participants ---

// createAdmin ...
func (ss *StorageSrv) createAdmin(nickName, web3Address string) error {
	// try to find admin by address
	var admin *Participant
	ss.psqlDB.Where("web3_address = ?", web3Address).Find(&admin)
	if admin.ID != "" {
		return nil
	}
	return ss.psqlDB.Create(&Participant{
		ID:          uuid.New().String(),
		NickName:    nickName,
		Web3Address: web3Address,
		Role:        AdminRole,
	}).Error
}

// CreateParticipant ...
func (ss *StorageSrv) CreateParticipant(nickName, web3Address string) (*Participant, error) {
	participant := &Participant{
		ID:          uuid.New().String(),
		NickName:    nickName,
		Web3Address: web3Address,
		Role:        ReaderRole,
	}
	//
	if err := ss.psqlDB.Create(participant).Error; err != nil {
		// todo error.As()
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, ErrParticipantAlreadyExists
		}

		return nil, err
	}

	return participant, nil
}

// FindParticipant ...
func (ss *StorageSrv) FindParticipant(nickName, emailAddress, web3Address string) (participant *Participant) {
	err := ss.psqlDB.Where("nick_name = ? OR email_address = ? OR web3_address = ?",
		nickName, emailAddress, web3Address).Find(&participant).Error
	if err != nil {
		ss.log.Error(err.Error())
		return nil
	}

	if participant.ID != "" {
		return participant
	}

	return nil
}

func (ss *StorageSrv) GetParticipantByAddress(address string) (participant *Participant, err error) {
	if err = ss.psqlDB.Where("web3_address = ?", address).Find(&participant).Error; err != nil {
		return nil, err
	}

	if participant.ID == "" {
		err = ErrParticipantNotExists
	}

	return
}

func (ss *StorageSrv) GetParticipantWorkByID(workID string) (work *ParticipantsWork, err error) {
	if ss.psqlDB.Where("work_id = ?", workID).Find(&work).Error != nil {
		return nil, ss.psqlDB.Where("work_id = ?", workID).Find(&work).Error
	}

	if work.ID == "" {
		err = ErrWorkNotExists
	}

	return
}

// func (ss *StorageSrv) GetAuthor(address string) (participant *Participant) {
// 	if err := ss.psqlDB.Where("web3_address = ?", address).Find(&participant).Error; err != nil {
// 		ss.log.Error(err.Error())
// 		return nil
// 	}

// 	if participant.ID == "" {
// 		return nil
// 	}
// 	return
// }

func (ss *StorageSrv) GetParticipantById(id string) (participant *Participant) {
	if err := ss.psqlDB.Where("id = ?", id).Find(&participant).Error; err != nil {
		ss.log.Error(err.Error())
		return nil
	}

	if participant.ID == "" {
		return nil
	}
	return
}

func (ss *StorageSrv) createWorkOfParticipant(authorID, workID string) error {
	return ss.psqlDB.Create(&ParticipantsWork{
		ID:            uuid.New().String(),
		ParticipantID: authorID,
		WorkID:        workID,
		Status:        ReviewWorkStatus,
		CreatedAt:     time.Now().UTC(),
	}).Error
}

func (ss *StorageSrv) getPendingWorks() []*ParticipantsWork {
	var works []*ParticipantsWork
	if err := ss.psqlDB.Where("status = ?", PreReviewWorkStatus).Find(&works).Error; err != nil {
		ss.log.Error(err.Error())
		return nil
	}
	if len(works) == 0 {
		return nil
	}
	return works
}

// updatewWorkStatus ...
func (ss *StorageSrv) updatewWorkStatus(workID string, newStatus WorkStatus) error {
	return ss.psqlDB.Model(ParticipantsWork{}).Where("work_id = ?", workID).
		Update("status", newStatus).Error
}

func (ss *StorageSrv) getParticipantIDOrNil(participantAddress string) string {
	var participant *Participant
	if err := ss.psqlDB.Where("web3_address = ?", participantAddress).Find(&participant).Error; err != nil {
		ss.log.Error(fmt.Sprintf("while getParticipantIDOrNil: %v", err))
		return ""
	}
	return participant.ID
}

func (ss *StorageSrv) getParticipantAddressOrNil(participantID string) string {
	var participant *Participant
	if err := ss.psqlDB.Where("id = ?", participantID).Find(&participant).Error; err != nil {
		ss.log.Error(fmt.Sprintf("while getParticipantIDOrNil: %v", err))
		return ""
	}
	return participant.Web3Address
}

func (ss *StorageSrv) getParticipantWorks(participantID string) []*ParticipantsWork {
	var works []*ParticipantsWork
	if err := ss.psqlDB.Where("participant_id = ?", participantID).Find(&works).Error; err != nil {
		ss.log.Error(fmt.Sprintf("while getParticipantWorks, err: %v", err))
		return nil
	}
	if len(works) == 0 {
		return nil
	}
	return works
}

func (ss *StorageSrv) GetAllParticipantsWorks() []*ParticipantsWork {
	var works []*ParticipantsWork
	if err := ss.psqlDB.Where("status <> ?", DeclinedWorkStatus).Find(&works).Error; err != nil {
		ss.log.Error(fmt.Sprintf("while getAllParticipantWorks, err: %v", err))

		return nil
	}

	return works
}

func (ss *StorageSrv) removeParticipantsWorkByID(workID string) error {
	return ss.psqlDB.Where("work_id = ?", workID).Delete(&ParticipantsWork{}).Error
}

/// Purchase works

func (ss *StorageSrv) PurchaseWork(participantID, workID string) error {
	purpose := &ParticipantsPurpose{
		ID:            uuid.New().String(),
		ParticipantID: participantID,
		WorkID:        workID,
	}
	if err := ss.psqlDB.Create(purpose).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (ss *StorageSrv) GetAvailableWorks(participantID string) []*ParticipantsPurpose {
	var works []*ParticipantsPurpose
	if err := ss.psqlDB.Where("participant_id = ?", participantID).Find(&works).Error; err != nil {
		ss.log.Error(fmt.Sprintf("while gtAvailableWorks, err: %v", err))
		return nil
	}
	if len(works) == 0 {
		return nil
	}

	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	for _, work := range works {
		if work.CreatedAt.UTC().Unix() < twoDaysAgo.Unix() {
			continue
		}
		fmt.Println("HERE!!!")
	}

	return works
}

func (ss *StorageSrv) PurchasedWorkOrNot(participantID, workID string) bool {
	var work *ParticipantsPurpose
	if err := ss.psqlDB.Where("participant_id = ? and work_id = ?", participantID, workID).First(&work).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		ss.log.Error(fmt.Sprintf("while PurchasedWorkOrNot, err: %v", err))
		return false
	}
	if work.ID == "" {
		return false
	}

	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	return work.CreatedAt.UTC().Unix() >= twoDaysAgo.Unix()
}

func (ss *StorageSrv) GetPurchasedByParticipantID(participantID string) (workIDs []string) {
	var works []*ParticipantsPurpose
	if err := ss.psqlDB.Where("participant_id = ?", participantID).Find(&works).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		ss.log.Error(fmt.Sprintf("while PurchasedWorkOrNot, err: %v", err))
		return
	}

	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	for _, work := range works {
		if work.CreatedAt.UTC().Unix() >= twoDaysAgo.Unix() {
			workIDs = append(workIDs, work.WorkID)
		}
	}
	return
}

func (ss *StorageSrv) removeWorkFromPurposes(workID string) error {
	return ss.psqlDB.Where("work_id = ?", workID).Delete(&ParticipantsPurpose{}).Error
}

func (ss *StorageSrv) CreateBookmark(participantID, workID string) error {
	// create a new bookmark for the participant
	if err := ss.psqlDB.Create(&ParticipantsBookmark{
		ID:            uuid.NewString(),
		ParticipantID: participantID,
		WorkID:        workID,
	}).Error; err != nil {
		return err
	}
	return nil
}

// TODO add log
func (ss *StorageSrv) RemoveBookmark(participantID, workID string) error {
	if err := ss.psqlDB.Where("participant_id = ? AND work_id = ?", participantID, workID).
		Delete(&ParticipantsBookmark{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (ss *StorageSrv) removeWorkFromBookmarks(workID string) error {
	if err := ss.psqlDB.Where("work_id = ?", workID).Delete(&ParticipantsBookmark{}).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (ss *StorageSrv) getPacticipantBookmarkIDs(participantID string) (workIDs []string, err error) {
	workIDs = []string{}
	if err := ss.psqlDB.
		Model(ParticipantsBookmark{}).
		Select("work_id").
		Where("participant_id = ?", participantID).
		Scan(&workIDs).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return workIDs, nil
}
