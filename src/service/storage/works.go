package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PutWork returns workID(uuid) or error
func (ss *StorageSrv) PutWork(ctx context.Context, work *Work) (string, error) {
	work.CreatedAt = time.Now().UTC()
	work.ID = uuid.New().String()
	collection := ss.mongoDB.Collection(collectionWorks)
	if collection == nil {
		panic(fmt.Errorf("works collection is nil"))
	}

	if _, err := collection.InsertOne(ctx, work); err != nil {
		return "", err
	}
	return work.ID, nil
}

// Returns the certain work by id
func (ss *StorageSrv) GetWorkByID(workID string) (response *WorkResponse, err error) {
	// postgres work
	participantsWork, err := ss.GetParticipantWorkByID(workID)
	if err != nil {
		if errors.Is(err, ErrWorkNotExists) {
			err = nil
		}

		return
	}

	// mongo work
	filter := make(map[string]interface{})
	filter["id"] = participantsWork.WorkID
	mongoWork, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if len(mongoWork) > 1 {
		ss.log.Error(fmt.Sprintf("the length of response(%d) is more than 1 work", len(mongoWork)))
		return nil, fmt.Errorf("something went wrong")
	}

	// get author info
	author, err := ss.GetAuthorById(participantsWork.ParticipantID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("while getting the inforamationa about author with id %s, err: %v",
			participantsWork.ParticipantID, err))
	}
	participant := ss.GetParticipantById(participantsWork.ParticipantID)

	// TODO
	return ss.buildWorkResponse(mongoWork[0], author, participant, true, false), nil
}

// Returns all works
func (ss *StorageSrv) GetAllWorks(readerAddress string) (response []*WorkResponse, err error) {
	var readerID string
	if readerAddress != "" {
		readerID = ss.getParticipantIDOrNil(readerAddress)
		if readerID == "" {
			return nil, fmt.Errorf("haven't got the reader: %s", readerAddress)
		}
	}

	participantsWorks := ss.getAllParticipantsWorks()
	response = make([]*WorkResponse, 0, len(participantsWorks))
	if len(participantsWorks) == 0 {
		return
	}

	filter := make(map[string]interface{})
	var worksIds []string
	for _, work := range participantsWorks {
		worksIds = append(worksIds, work.WorkID)
	}
	filter["id"] = worksIds

	// gonna search by ids
	mongoWorks, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	participant := ss.GetParticipantById(readerID)

	for _, work := range participantsWorks {
		for _, mWork := range mongoWorks {
			if work.WorkID == mWork.ID {
				// get author info
				authorInfo, err := ss.GetAuthorById(work.ParticipantID)
				if err != nil {
					ss.log.Error(fmt.Sprintf("while getting the inforamationa about author with id %s, err: %v", work.ParticipantID, err))
				}
				authorBasicInfo := ss.GetParticipantById(work.ParticipantID)
				// the participant status is Reader just to show the annotation and
				// other preview information of work
				purchased := false
				if ss.PurchasedWorkOrNot(readerID, work.WorkID) || readerID == authorBasicInfo.ID {
					purchased = true
				}
				// verify whether the work is open or the reader is authorof this work/validator/admin
				work, content := work.IsShow(participant, purchased)
				if work {
					response = append(response,
						ss.buildWorkResponse(mWork, authorInfo, authorBasicInfo, content, ss.BookmarkedWorkOrNot(readerID, mWork.ID)),
					)
				}
			}
		}
	}

	return response, nil
}

// GetWorkByFilter ...
func (ss *StorageSrv) GetWorkByKeyWords(readerAddress string, keyWords []string) (response []*WorkResponse, err error) {
	var readerID string
	if readerAddress != "" {
		readerID = ss.getParticipantIDOrNil(readerAddress)
		if readerID == "" {
			return nil, fmt.Errorf("haven't got the reader: %s", readerAddress)
		}
	}

	participantsWorks := ss.getAllParticipantsWorks()
	if participantsWorks == nil {
		return nil, fmt.Errorf("there are no works in the library:(")
	}

	// gonna search by ids
	mongoWorks, err := ss.getWorksByKeyWords(keyWords)
	if err != nil {
		return nil, err
	}

	for _, work := range participantsWorks {
		for _, mWork := range mongoWorks {
			if work.WorkID == mWork.ID {
				// get author info
				author, err := ss.GetAuthorById(work.ParticipantID)
				if err != nil {
					ss.log.Error(fmt.Sprintf("while getting the inforamationa about author with id %s, err: %v", work.ParticipantID, err))
				}
				participant := ss.GetParticipantById(work.ParticipantID)
				// the participant status is Reader just to show the annotation and
				// other preview information of work
				purchased := false
				if ss.PurchasedWorkOrNot(readerID, work.ID) || readerID == author.ID {
					purchased = true
				}
				response = append(response,
					ss.buildWorkResponse(mWork, author, participant, purchased, ss.BookmarkedWorkOrNot(readerID, mWork.ID)),
				)
			}
		}
	}
	return response, nil
}

// GetPurchasedWorks ...
func (ss *StorageSrv) GetPurchasedWorks(ctx context.Context, readerAddress string) (response []*WorkResponse, err error) {
	readerID := ss.getParticipantIDOrNil(readerAddress)
	if readerID == "" {
		return nil, fmt.Errorf("haven't got the reader: %s", readerAddress)
	}

	workIDs := ss.GetPurchasedByParticipantID(readerID)
	if len(workIDs) == 0 {
		return []*WorkResponse{}, nil
	}

	filter := make(map[string]interface{})
	filter["id"] = workIDs
	mongoWorks, err := ss.getWorksByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, mWork := range mongoWorks {
		authorBasicInfo := ss.GetParticipantById(mWork.AuthorID)
		// get author info
		author, err := ss.GetAuthorById(mWork.AuthorID)
		if err != nil {
			ss.log.Error(fmt.Sprintf("while getting the inforamationa about author with id %s, err: %v", mWork.AuthorID, err))
		}
		response = append(response,
			ss.buildWorkResponse(mWork, author, authorBasicInfo, true, ss.BookmarkedWorkOrNot(readerID, mWork.ID)),
		)
	}
	return response, nil
}

func (ss *StorageSrv) GetWorksByAuthorAddress(
	readerAddress,
	authorAddress string,
) (response []*WorkResponse, err error) {
	readerID := ss.getParticipantIDOrNil(readerAddress)
	if readerID == "" {
		return nil, fmt.Errorf("there is no participant with the address %s", readerAddress)
	}

	participant, err := ss.GetParticipantByAddress(authorAddress)
	if err != nil {
		return nil, err
	}

	participantWorks := ss.getParticipantWorks(participant.ID)
	if participantWorks == nil {
		return nil, fmt.Errorf("there is no works for participant with the address %s", authorAddress)
	}

	filter := make(map[string]interface{})
	for _, work := range participantWorks {
		filter["id"] = work.ID
	}
	// get mongo works
	works, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// get the author information
	// according to the request `get by author id` it means the single author
	// for all works
	author, err := ss.GetAuthorById(participantWorks[0].ParticipantID)
	if err != nil {
		return nil, err
	}

	for _, partWork := range participantWorks {
		for _, work := range works {
			if partWork.WorkID == work.ID {
				purchased := false
				if ss.PurchasedWorkOrNot(readerID, work.ID) || readerID == author.ID {
					purchased = true
				}
				response = append(response, ss.buildWorkResponse(work, author, participant, purchased, ss.BookmarkedWorkOrNot(readerID, readerID)))
			}
		}
	}

	return response, nil
}

// Returns all works in
func (ss *StorageSrv) GetWorksByAuthorID(authorID string) ([]*Work, error) {
	// gonna search by author_id key
	return ss.getWorksByFilter(context.TODO(), map[string]interface{}{"author_id": authorID})
}

// Returns all works in
func (ss *StorageSrv) GetPendingWorks() (response []*WorkResponse, err error) {
	participantsWorks := ss.getPendingWorks()
	if participantsWorks == nil {
		return nil, nil
	}

	filter := make(map[string]interface{})
	var worksIds []string
	for _, work := range participantsWorks {
		worksIds = append(worksIds, work.WorkID)
	}
	filter["id"] = worksIds

	// gonna search by ids
	mongoWorks, err := ss.getWorksByFilter(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for _, work := range participantsWorks {
		for _, mWork := range mongoWorks {
			if work.WorkID == mWork.ID {
				// get author info for this work
				author, err := ss.GetAuthorById(work.ParticipantID)
				if err != nil {
					ss.log.Error(fmt.Sprintf("while getting the inforamationa about author with id %s, err: %v", work.ParticipantID, err))
				}
				participant := ss.GetParticipantById(work.ParticipantID)
				response = append(response, ss.buildWorkResponse(mWork, author, participant, true, false))
			}
		}
	}
	return
}

func (ss *StorageSrv) getWorksByFilter(ctx context.Context, options map[string]interface{}) (works []*Work, err error) {
	// pack all filter opt together
	var filter interface{}
	for key, value := range options {
		switch values := value.(type) {
		case string:
			filter = bson.M{key: values}
		case []string:
			filter = bson.M{key: bson.M{"$in": values}}

			// oids := make([]bson, len(values))

			// for _, v := range values {
			// filter = append(filter, bson.E{Key: "id", Value: values})
			// 	oids[i] = bson.ObjectIdHex(ids[i])

			// 	query := bson.M{"_id": bson.M{"$in": oids}}
			// 	filter = append(filter, bson.E{Key: key, Value: v})
			// }
		default:
			// filter = append(filter, bson.E{Key: key, Value: value})
		}
	}

	collection := ss.mongoDB.Collection(collectionWorks)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	// make a request with the filter
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ss.log.Error(fmt.Sprintf("while finding works, err: %v", err))
		return nil, err
	}

	if err = cur.All(ctx, &works); err != nil {
		ss.log.Error(fmt.Sprintf("while decoding works, err: %v", err))
		return nil, err
	}

	return works, nil
}

func (ss *StorageSrv) getWorksByKeyWords(keyWords []string) (works []*Work, err error) {
	keyWords = removeDuplicate(keyWords)
	collection := ss.mongoDB.Collection(collectionWorks)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	for _, key := range keyWords {
		filter := bson.D{{"$text", bson.D{{"$search", fmt.Sprintf("\"%s\"", key)}}}}
		// make a request with the filter
		cur, err := collection.Find(context.Background(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			ss.log.Error(fmt.Sprintf("while finding works, err: %v", err))
			return nil, err
		}
		var filterWorks []*Work
		if err = cur.All(context.Background(), &filterWorks); err != nil {
			ss.log.Error(fmt.Sprintf("while decoding works, err: %v", err))
			return nil, err
		}
		works = append(works, filterWorks...)
	}

	return removeWorksDuplicate(works), nil
}

func (ss *StorageSrv) removeWorkByID(workID string) error {
	// pack all filter opt together
	filter := bson.M{"id": bson.M{"$in": []string{workID}}}
	collection := ss.mongoDB.Collection(collectionWorks)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	// make a request with the filter
	if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		ss.log.Error(fmt.Sprintf("while finding works, err: %v", err))
		return err
	}

	return nil
}

func removeWorksDuplicate(works []*Work) []*Work {
	allKeys := make(map[string]bool)
	result := []*Work{}
	for _, work := range works {
		if _, value := allKeys[work.ID]; !value {
			allKeys[work.ID] = true
			result = append(result, work)
		}
	}
	return result
}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
