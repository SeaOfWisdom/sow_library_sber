package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionWorks = "works"

	collectionAuthors = "authors"

	collectionValidators = "validators"

	collectionWorkReviews = "work_reviews"
)

func addIndexOnWorks(works *mongo.Collection) {
	defaultLang, langOverride := "none", "none"
	model := mongo.IndexModel{Keys: bson.D{{"annotation", "text"}, {"name", "text"}},
		Options: &options.IndexOptions{DefaultLanguage: &defaultLang, LanguageOverride: &langOverride}}
	if _, err := works.Indexes().CreateOne(context.Background(), model); err != nil {
		panic(err)
	}
}

// ApproveWork ...
func (ss *StorageSrv) ApproveWork(workID string) error {
	// get work by id
	works, err := ss.getWorksByFilter(context.TODO(), map[string]interface{}{"work_id": workID})
	if err != nil {
		return err
	}

	if len(works) != 1 {
		return fmt.Errorf("wrong len of returned works(%d)", len(works))
	}
	work := works[0]

	if err := ss.updatewWorkStatus(work.ID, ReviewWorkStatus); err != nil {
		return err
	}

	return nil
}

// RemoveWork ...
func (ss *StorageSrv) RemoveWork(workID string) error {
	participantsWork, err := ss.GetParticipantWorkByID(workID)
	if err != nil {
		ss.log.Error(fmt.Sprint("RemoveWork: error get participant by work id, err: %v", err))

		return err
	}

	filter := make(map[string]interface{})
	filter["id"] = participantsWork.WorkID
	mongoWork, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return err
	}

	if len(mongoWork) > 1 {
		ss.log.Error(fmt.Sprintf("the length of response(%d) is more than 1 work", len(mongoWork)))
		return fmt.Errorf("something went wrong")
	}

	// remove from MongoDB
	if err := ss.removeWorkByID(workID); err != nil {
		ss.log.Error(fmt.Sprintf("while removing the work from MongoDB, err: %v", err))
		return fmt.Errorf("something went wrong")
	}

	// remove from PostgreSQL
	if err := ss.removeParticipantsWorkByID(workID); err != nil {
		ss.log.Error(fmt.Sprintf("while removing the work from PostgreSQL(participantsWorks), err: %v", err))
		return fmt.Errorf("something went wrong")
	}

	// remove from Bookmarks
	if err := ss.removeWorkFromBookmarks(workID); err != nil {
		ss.log.Error(fmt.Sprintf("while removing the work from PostgreSQL(bookmarks), err: %v", err))
		return fmt.Errorf("something went wrong")
	}

	// remove from Purposes
	if err := ss.removeWorkFromPurposes(workID); err != nil {
		ss.log.Error(fmt.Sprintf("while removing the work from PostgreSQL(purposes), err: %v", err))
		return fmt.Errorf("something went wrong")
	}

	return nil
}

// ConfirmWork ...
func (ss *StorageSrv) ConfirmWork(ctx context.Context, workID string) error {
	// get work by id
	works, err := ss.getWorksByFilter(ctx, map[string]interface{}{"id": workID})
	if err != nil {
		return err
	}

	if len(works) != 1 {
		return fmt.Errorf("wrong len of returned works(%d)", len(works))
	}
	work := works[0]

	if err := ss.updatewWorkStatus(work.ID, OpenWorkStatus); err != nil {
		return err
	}

	return nil
}

func (ss *StorageSrv) DeclineWork(ctx context.Context, workID string) error {
	// get work by id
	works, err := ss.getWorksByFilter(ctx, map[string]interface{}{"id": workID})
	if err != nil {
		return err
	}

	if len(works) != 1 {
		return fmt.Errorf("wrong len of returned works(%d)", len(works))
	}
	work := works[0]

	if err := ss.updatewWorkStatus(work.ID, DeclinedWorkStatus); err != nil {
		return err
	}

	return nil
}

// TODO
func (ss *StorageSrv) UpdateWork(work *Work) error {
	work.UpdatedAt = time.Now().UTC()
	// if _, err := ss.worksCollection.UpdateOne(context.Background(), work); err != nil {
	// 	return err
	// }
	return nil
}
