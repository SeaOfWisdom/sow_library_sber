package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ss *StorageSrv) CreateParticipantsWorkReview(validatorID, workID string) (review *ParticipantsWorkReview, err error) {
	ss.psqlDB.Where("participant_id = ? AND work_id = ?", validatorID, workID).Find(&review)
	if review.ID != "" {
		return review, nil
	}
	review = &ParticipantsWorkReview{
		ID:            uuid.New().String(),
		ParticipantID: validatorID,
		WorkID:        workID,
		Status:        WorkReviewInProgress,
	}
	if err := ss.psqlDB.Create(review).Error; err != nil {
		return nil, err
	}
	return
}

func (ss *StorageSrv) FindParticipantsWorkReviewByValidator(validatorID, workID string) (review *ParticipantsWorkReview, err error) {
	if err := ss.psqlDB.Where("participant_id = ? AND work_id = ?",
		validatorID, workID).Find(&review).Error; err != nil {
		return nil, err
	}
	return
}

func (ss *StorageSrv) FindParticipantsWorkReviews(workID string) (review []*ParticipantsWorkReview, err error) {
	if err := ss.psqlDB.Where("work_id = ?", workID).Find(&review).Error; err != nil {
		return nil, err
	}
	return
}

func (ss *StorageSrv) UpdateParticipantsWorkReviewStatus(id string, newStatus WorkReviewStatus) error {
	fmt.Println("ID: ", id)
	return ss.psqlDB.Model(ParticipantsWorkReview{}).Where("id = ?", id).
		Update("status", newStatus).Error
}

func (ss *StorageSrv) UpdateOrCreateWorkReview(ctx context.Context, validatorID string, review *WorkReview) (*WorkReview, error) {
	collection := ss.mongoDB.Collection(collectionWorkReviews)
	if collection == nil {
		panic(fmt.Errorf("work_reviews collection is nil"))
	}

	participantsReview, err := ss.CreateParticipantsWorkReview(validatorID, review.WorkID)
	if err != nil {
		return nil, nil
	}

	currentReview, err := ss.GetWorkReviewByID(participantsReview.ID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("while getting the workReview by validator and author id, err: %s", err))
		return nil, err
	}

	if currentReview == nil {
		review.ID = participantsReview.ID
		review.CreatedAt = time.Now().UTC()
		if _, err := collection.InsertOne(ctx, review); err != nil {
			return nil, err
		}

		return review, nil
	}

	if currentReview.Body == nil {
		currentReview.Body = &WorkReviewBody{}
	}

	if review.Body.Questionnaire != nil {
		if review.Body.Questionnaire.Questions != nil {
			if currentReview.Body.Questionnaire == nil {
				currentReview.Body.Questionnaire = &WorkReviewQuestionnaire{}
			}
			currentReview.Body.Questionnaire.Questions = review.Body.Questionnaire.Questions
		}
	}
	if review.Body.Review != "" {
		currentReview.Body.Review = review.Body.Review
	}

	if err := ss.UpdateWorkReview(context.TODO(), currentReview); err != nil {
		return nil, nil
	}
	return currentReview, nil
}

// SubmitWorkReview ...
func (ss *StorageSrv) SubmitWorkReview(participantReview *ParticipantsWorkReview) error {
	if err := ss.UpdateParticipantsWorkReviewStatus(participantReview.ID, participantReview.Status); err != nil {
		return err
	}

	currentReview, err := ss.GetWorkReviewByID(participantReview.ID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("while getting the workReview by its id, err: %s", err))

		return err
	}

	return ss.UpdateWorkReview(context.TODO(), currentReview)
}

// SubmitWorkReview ...
func (ss *StorageSrv) GetWorkReviewByID(id string) (review *WorkReview, err error) {
	filter := bson.M{"id": id}
	collection := ss.mongoDB.Collection(collectionWorkReviews)
	if collection == nil {
		panic(fmt.Errorf("work_reviews collection is nil"))
	}

	// make a request with the filter
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ss.log.Error(fmt.Sprintf("while finding authors, err: %v", err))
		return nil, err
	}

	for cur.Next(context.Background()) {
		if err := cur.Decode(&review); err != nil {
			ss.log.Error(fmt.Sprintf("Error decoding document, err: %v", err))
			return nil, err
		}
		return review, nil
	}
	return
}

// TODO rewrite without for loop
func (ss *StorageSrv) GetWorkReviewByWorkId(workID string) (review *WorkReview, err error) {
	filter := bson.M{"work_id": workID}
	collection := ss.mongoDB.Collection(collectionWorkReviews)
	if collection == nil {
		panic(fmt.Errorf("work_reviews collection is nil"))
	}

	// make a request with the filter
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ss.log.Error(fmt.Sprintf("while finding authors, err: %v", err))
		return nil, err
	}

	for cur.Next(context.Background()) {
		if err := cur.Decode(&review); err != nil {
			ss.log.Error(fmt.Sprintf("Error decoding document, err: %v", err))
			return nil, err
		}
		return review, nil
	}
	return
}

func (ss *StorageSrv) GetReviewByValidatorAndWorkID(validatorID, workID string) (review *WorkReview, err error) {
	participantReview, err := ss.FindParticipantsWorkReviewByValidator(validatorID, workID)
	if err != nil {
		return nil, err
	}
	if participantReview == nil {
		return nil, nil
	}

	return ss.GetWorkReviewByID(participantReview.ID)
}

func (ss *StorageSrv) GetReviewByAuthorAndWorkID(authorID, workID string) (reviews []*WorkReview, err error) {
	participantReviews, err := ss.FindParticipantsWorkReviews(workID)
	if err != nil {
		return nil, err
	}
	if len(participantReviews) == 0 {
		return nil, nil
	}

	for _, participantReview := range participantReviews {
		review, err := ss.GetWorkReviewByID(participantReview.ID)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return
}

func (ss *StorageSrv) GetReviewsByWorkId(workID string) (reviews []*ParticipantsWorkReview, err error) {
	return ss.FindParticipantsWorkReviews(workID)
}

// func (ss *StorageSrv) geWorkReByFilter(options map[string]interface{}, preRead bool) (works []*Work, err error) {
// 	// pack all filter opt together
// 	filter := bson.D{}
// 	for key, value := range options {
// 		filter = append(filter, bson.E{Key: key, Value: value})
// 	}

// 	collection := ss.mongoDB.Collection(collectionWorkReviews)
// 	if collection == nil {
// 		panic(fmt.Errorf("work_reviews collection is nil"))
// 	}

// 	// make a request with the filter
// 	cur, err := collection.Find(context.Background(), filter)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil
// 		}
// 		ss.log.Error(fmt.Sprintf("while finding works, err: %v", err))
// 		return nil, err
// 	}

// 	if err = cur.All(context.Background(), &works); err != nil {
// 		ss.log.Error(fmt.Sprintf("while decoding works, err: %v", err))
// 		return nil, err
// 	}

// 	return works, nil
// }

func (ss *StorageSrv) UpdateWorkReview(ctx context.Context, review *WorkReview) error {
	review.UpdatedAt = time.Now().UTC()
	filter := bson.M{"id": review.ID}
	collection := ss.mongoDB.Collection(collectionWorkReviews)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	update := bson.M{
		"$set": review,
	}

	return collection.FindOneAndUpdate(ctx, filter, update).Err()
}
