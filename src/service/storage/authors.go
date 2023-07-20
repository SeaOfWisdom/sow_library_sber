package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ss *StorageSrv) CreateAuthor(postgesqlID, emailAddress, name, surname string) error {
	collection := ss.mongoDB.Collection(collectionAuthors)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	author := Author{
		ID:           postgesqlID,
		Name:         name,
		Surname:      surname,
		EmailAddress: emailAddress,
		CreatedAt:    time.Now().UTC(),
	}

	if _, err := collection.InsertOne(context.Background(), author); err != nil {
		return err
	}

	return nil
}

// TODO rewrite without for loop
func (ss *StorageSrv) GetAuthorById(ctx context.Context, id string) (author *Author, err error) {
	filter := bson.M{"id": id}
	collection := ss.mongoDB.Collection(collectionAuthors)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	// make a request with the filter
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}

		return nil, err
	}

	for cur.Next(context.Background()) {
		err = cur.Decode(&author)

		return
	}

	return
}

func (ss *StorageSrv) getAuthorsByFilter(ctx context.Context, options map[string]interface{}, preRead bool) (works []*Work, err error) {
	// pack all filter opt together
	filter := bson.D{}
	for key, value := range options {
		filter = append(filter, bson.E{Key: key, Value: value})
	}

	collection := ss.mongoDB.Collection(collectionAuthors)
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
		ss.log.Errorf("while decoding works, err: %v", err)

		return nil, err
	}

	return works, nil
}

func (ss *StorageSrv) UpdateAuthorInfo(ctx context.Context, author *Author) error {
	author.UpdatedAt = time.Now().UTC()
	filter := bson.M{"id": author.ID}
	collection := ss.mongoDB.Collection(collectionAuthors)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	update := bson.M{
		"$set": author,
	}

	return collection.FindOneAndUpdate(ctx, filter, update).Err()
}
