package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/SeaOfWisdom/sow_library/src/config"

	"gorm.io/gorm"
)

// StorageSrv ...
type StorageSrv struct {
	/* log */
	log *zap.Logger
	/* PostgreSQL */
	psqlDB *gorm.DB
	/* MongoDB */
	mongoDB *mongo.Database

	// collections
	//	worksCollection *mongo.Collection
}

// NewStorageSrv ...
func NewStorageSrv(cfg *config.Config, postresDB *gorm.DB, mongoDB *mongo.Database) *StorageSrv {
	// verify the existing of all collections
	collection := mongoDB.Collection(collectionWorks)
	if collection == nil {
		panic(fmt.Errorf("works collection is nil"))
	}
	addIndexOnWorks(collection)

	collection = mongoDB.Collection(collectionAuthors)
	if collection == nil {
		panic(fmt.Errorf("authors collection is nil"))
	}

	collection = mongoDB.Collection(collectionValidators)
	if collection == nil {
		panic(fmt.Errorf("validators collection is nil"))
	}

	collection = mongoDB.Collection(collectionWorkReviews)
	if collection == nil {
		panic(fmt.Errorf("work_reviews collection is nil"))
	}

	ss := &StorageSrv{
		log:     zap.NewExample(),
		psqlDB:  postresDB,
		mongoDB: mongoDB,
	}
	// go postgreSQL migrations
	if err := ss.psqlDB.AutoMigrate(Participant{}); err != nil {
		panic(err)
	}
	if err := ss.psqlDB.AutoMigrate(ParticipantsWork{}); err != nil {
		panic(err)
	}
	if err := ss.psqlDB.AutoMigrate(ParticipantsPurpose{}); err != nil {
		panic(err)
	}
	if err := ss.psqlDB.AutoMigrate(ParticipantsBookmark{}); err != nil {
		panic(err)
	}
	if err := ss.psqlDB.AutoMigrate(ParticipantsWorkReview{}); err != nil {
		panic(err)
	}
	// create admins from the config if they don't exist
	for nickName, address := range config.AdminAddresses {
		if err := ss.createAdmin(nickName, address); err != nil {
			panic(err)
		}
	}
	return ss
}

func (ss *StorageSrv) CreateWork(ctx context.Context, authorID string, work *Work) (string, error) {
	work.AuthorID = authorID
	workID, err := ss.PutWork(ctx, work)
	if err != nil {
		return "", err
	}

	if err := ss.createWorkOfParticipant(authorID, workID); err != nil {
		return "", err
	}

	return workID, nil
}

// the participant status is Reader just to show the annotation and
// other preview information of work
func (ss *StorageSrv) buildWorkResponse(
	//	role ParticipantRole,
	work *Work,
	author *Author,
	participant *Participant,
	showContent bool,
	bookmarked bool,
) (workResp *WorkResponse) {
	if author != nil {
		work.AuthorID = author.ID
	}
	work.Price = "50000000000000000000" // wei
	workResp = new(WorkResponse)
	workResp.Work = work
	workResp.Author = &AuthorResponse{
		BasicInfo:  &Participant{},
		AuthorInfo: author,
	}
	if participant != nil {
		workResp.Author.BasicInfo = participant
	}

	// content is available in cases:

	if !showContent {
		workResp.Work.Content = nil
	}

	workResp.Bookmarked = bookmarked

	// switch status {
	// case OpenWorkStatus:
	// 	return
	// case ReviewWorkStatus:
	// 	if role == AdminRole {
	// 		return
	// 	}
	// 	workResp.Work.Content = nil
	// case PreReviewWorkStatus:
	// case DeclinedWorkStatus:
	// 	if role == AdminRole {
	// 		return
	// 	}
	// 	workResp.Work.Annotation = ""
	// 	workResp.Work.Description = ""
	// 	workResp.Work.Content = nil
	// }
	return
}
