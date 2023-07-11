package storage

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (ss *StorageSrv) GetBookmarksByParticipantID(participantID string) (bookmarks []*WorkResponse, err error) {
	// get the participant's id
	workIDs, err := ss.getPacticipantBookmarkIDs(participantID)
	if err != nil {
		return
	}
	if len(workIDs) == 0 {
		return nil, nil
	}

	// get works by the ids
	filter := make(map[string]interface{})
	filter["id"] = workIDs
	// gonna search by ids
	mongoWorks, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for _, mWork := range mongoWorks {
		// get author info
		author, err := ss.GetAuthorById(mWork.AuthorID)
		if err != nil {
			ss.log.Error(fmt.Sprintf("while getting the inforamation regarding the author with id %s, err: %v", mWork.AuthorID, err))
		}
		participant := ss.GetParticipantById(mWork.AuthorID)
		// the participant status is Reader just to show the annotation and
		// other preview information of work
		bookmarks = append(bookmarks, ss.buildWorkResponse(mWork, author, participant, ss.PurchasedWorkOrNot(participantID, mWork.ID), true))
	}

	return bookmarks, nil
}

func (ss *StorageSrv) BookmarkedWorkOrNot(participantID, workID string) bool {
	var bookmark *ParticipantsBookmark
	if err := ss.psqlDB.Where("participant_id = ? and work_id = ?", participantID, workID).First(&bookmark).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		ss.log.Error(fmt.Sprintf("while BookmarkedWorkOrNot, err: %v", err))
		return false
	}
	if bookmark.ID == "" {
		return false
	}

	return true
}
