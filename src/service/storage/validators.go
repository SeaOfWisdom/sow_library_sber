package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ss *StorageSrv) CreateValidator(postgesqlID, emailAddress, name, surname string) error {
	collection := ss.mongoDB.Collection(collectionValidators)
	if collection == nil {
		panic(fmt.Errorf("validators collection is nil"))
	}

	validator := Validator{
		ID:           postgesqlID,
		Name:         name,
		Surname:      surname,
		EmailAddress: emailAddress,
		CreatedAt:    time.Now().UTC(),
	}

	if _, err := collection.InsertOne(context.Background(), validator); err != nil {
		return err
	}

	return nil
}

// TODO rewrite without for loop
func (ss *StorageSrv) GetValidatorById(id string) (validator *Validator, err error) {
	filter := bson.M{"id": id}
	collection := ss.mongoDB.Collection(collectionValidators)
	if collection == nil {
		panic(fmt.Errorf("validators collection is nil"))
	}

	// make a request with the filter
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrParticipantNotExists
		}

		return
	}

	for cur.Next(context.Background()) {
		if err := cur.Decode(&validator); err != nil {
			return nil, err
		}

		return validator, nil
	}

	return nil, nil
}

func (ss *StorageSrv) UpdateValidatorInfo(validator *Validator) error {
	validator.UpdatedAt = time.Now().UTC()
	filter := bson.M{"id": validator.ID}
	collection := ss.mongoDB.Collection(collectionValidators)
	if collection == nil {
		panic(fmt.Errorf("validators collection is nil"))
	}

	update := bson.M{
		"$set": validator,
	}

	return collection.FindOneAndUpdate(context.Background(), filter, update).Err()
}
