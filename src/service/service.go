package srv

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SeaOfWisdom/sow_library/src/log"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"
	contractor "github.com/SeaOfWisdom/sow_proto/contractor-srv"
)

const (
	ServiceName = "sow_library"
	faucetCount = "50000000000000000000" // 50 ETH
)

var (
	ErrNoReviews            = errors.New("there are no reviews")
	ErrValidationNotAllowed = errors.New("validation now allowed")
)

type LibrarySrv struct {
	log     *log.Logger
	storage *storage.StorageSrv
	//works   []*storage.WorkResponse

	contractorSrv contractor.ContractorServiceClient
}

// create

func NewLibrarySrv(log *log.Logger, str *storage.StorageSrv, contractorSrv contractor.ContractorServiceClient) *LibrarySrv {
	return &LibrarySrv{
		log:           log,
		storage:       str,
		contractorSrv: contractorSrv,
	}
}

func (ls *LibrarySrv) Start() {
	ls.MigrateFromMongo()
	// // get all works from the library
	// works, err := lb.GetAllWorks(config.AdminAddresses["chillhacker"])
	// if err != nil {
	// 	if strings.Contains(err.Error(), "there are no works in the library") {
	// 		return
	// 	}
	// 	panic(err)
	// }
	// lb.works = works
}

// --- Handles

func (ls *LibrarySrv) CreateParticipant(ctx context.Context, nickname, web3Address string) (*storage.Participant, error) {
	participant, err := ls.storage.CreateParticipant(nickname, web3Address)
	if err != nil {
		ls.log.Errorf("CreateParticipant: error create participant, err: %v", err)

		return nil, err
	}

	out, err := ls.contractorSrv.AddParticipant(ctx, &contractor.AccountRequest{
		Address: web3Address,
	})
	if err != nil {
		ls.log.Errorf("CreateParticipant: error add participant to the chain, err: %v", err)

		return nil, err
	}

	if out.ErrorMsg != "" {
		err = fmt.Errorf("error add participant to the chain, err: %s", out.ErrorMsg)

		ls.log.Errorf("CreateParticipant: %v", err)

		return nil, err
	}

	ls.log.Infof("CreateParticipant: successfully added participant to chain, tx hash - %s", out.TxHash)

	return participant, nil
}

func (ls *LibrarySrv) UpdateParticipantNickName(web3Address, newNickName string) (*storage.Participant, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(web3Address)
	if err != nil {
		ls.log.Errorf("UpdateParticipantNickName: error get participant with address %s, err: %v", web3Address, err)

		return nil, err
	}

	if err = ls.storage.UpdateParticipantNickName(participant.ID, newNickName); err != nil {
		ls.log.Errorf("UpdateParticipantNickName: error update participant nickname, err: %v", err)

		return nil, storage.ErrSomethingWentWrong
	}

	participant.NickName = newNickName
	return participant, nil
}

func (ls *LibrarySrv) BecomeAuthor(web3Address, emailAddress, name, surname string) (*storage.Participant, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(web3Address)
	if err != nil {
		ls.log.Errorf("BecomeAuthor: error get participant with address %s, err: %v", web3Address, err)
		return nil, err
	}
	participant.Role = storage.AuthorRole

	if err = ls.storage.UpdateParticipantRole(participant.ID, storage.AuthorRole); err != nil {
		ls.log.Errorf("BecomeAuthor: error update participant role, err: %v", err)

		return nil, fmt.Errorf("while updating the participant's role, err: %w", err)
	}

	// create a new record for the author
	if err = ls.storage.CreateAuthor(participant.ID, emailAddress, name, surname); err != nil {
		ls.log.Errorf("BecomeAuthor: error create author, err: %v", err)

		return nil, fmt.Errorf("while creating an author, err: %w", err)
	}

	// let's add the new author into SowLibrary
	txHash, err := ls.contractorSrv.MakeAuthor(context.Background(), &contractor.AccountRequest{
		Address: web3Address,
	})
	if err != nil {
		ls.log.Errorf("BecomeAuthor: error create author, err: %v", err)

		return nil, fmt.Errorf("while creating an author, err: %w", err)
	}

	ls.log.Infof("the new author was added with tx: %s", txHash)

	return participant, nil
}

func (ls *LibrarySrv) GetAuthor(ctx context.Context, web3Address string) (*storage.AuthorResponse, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(web3Address)
	if err != nil {
		ls.log.Errorf("GetAuthor: error get participant with address %s, err: %v", web3Address, err)

		return nil, err
	}

	// verify his role
	if participant.Role < storage.AuthorRole {
		return nil, fmt.Errorf("the participant is not author")
	}

	// get author info
	author, err := ls.storage.GetAuthorById(ctx, participant.ID)
	if err != nil {
		ls.log.Errorf("GetAuthor: error get the author by id, err: %v", err)

		return nil, err
	}

	return &storage.AuthorResponse{
		BasicInfo:  participant,
		AuthorInfo: author,
	}, nil
}

func (ls *LibrarySrv) GetValidator(ctx context.Context, web3Address string) (*storage.ValidatorResponse, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(web3Address)
	if err != nil {
		ls.log.Errorf("GetValidator: error get participant with address %s, err: %v", web3Address, err)

		return nil, err
	}

	// verify his role
	if participant.Role < storage.ValidatorRole {
		return nil, fmt.Errorf("the participant is not validator")
	}

	// get validator info
	validator, err := ls.storage.GetValidatorById(ctx, participant.ID)
	if err != nil {
		ls.log.Errorf("GetValidator: error get the validator by id, err, err: %v", err)

		return nil, err
	}

	return &storage.ValidatorResponse{
		BasicInfo:     participant,
		ValidatorInfo: validator,
	}, nil
}

func (ls *LibrarySrv) UpdateAuthorInfo(
	ctx context.Context,
	web3Address,
	emailAddress,
	name,
	middlename,
	surname,
	orcid,
	scholarShipProfile,
	language string,
	sciences []string,
) (*storage.Participant, error) {
	// try to get the author record
	authorResp, err := ls.GetAuthor(ctx, web3Address)
	if err != nil {
		return nil, err
	}

	if emailAddress != "" {
		// it needs verification to be added here
		authorResp.AuthorInfo.EmailAddress = emailAddress
	}

	if name != "" {
		authorResp.AuthorInfo.Name = name
	}

	if middlename != "" {
		authorResp.AuthorInfo.MiddleName = middlename
	}

	if surname != "" {
		authorResp.AuthorInfo.Surname = surname
	}

	if orcid != "" {
		authorResp.AuthorInfo.Orcid = orcid
	}

	if scholarShipProfile != "" {
		authorResp.AuthorInfo.ScholarShipProfile = scholarShipProfile
	}

	if language != "" {
		authorResp.AuthorInfo.Language = language
	}

	if len(sciences) > 0 {
		authorResp.AuthorInfo.Sciences = sciences
	}

	// update the Auhtor info in the storage(MongoDB)
	if err = ls.storage.UpdateAuthorInfo(ctx, authorResp.AuthorInfo); err != nil {
		ls.log.Errorf("UpdateAuthorInfo: error update author info, err: %v", err)

		return nil, err
	}

	return authorResp.BasicInfo, nil
}

/// Validator

func (ls *LibrarySrv) BecomeValidator(ctx context.Context, web3Address, emailAddress, name, surname string) (*storage.Participant, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(web3Address)
	if err != nil {
		ls.log.Errorf("BecomeValidator: error get participant with address %s, err: %v", web3Address, err)

		return nil, err
	}

	participant.Role = storage.ValidatorRole

	if err = ls.storage.UpdateParticipantRole(participant.ID, storage.ValidatorRole); err != nil {
		ls.log.Errorf("BecomeValidator: error update participant role, err: %v", err)

		return nil, fmt.Errorf("while updating the participant's role, err: %s", err)
	}

	// create a new record for the validator
	if err := ls.storage.CreateValidator(ctx, participant.ID, emailAddress, name, surname); err != nil {
		ls.log.Errorf("BecomeValidator: error create validator, err: %v", err)

		return nil, fmt.Errorf("while creating validator, err: %v", err)
	}

	// let's add the new reviewer into SowLibrary
	txHash, err := ls.contractorSrv.MakeReviewer(ctx, &contractor.AccountRequest{
		Address: web3Address,
	})
	if err != nil {
		ls.log.Errorf("BecomeValidator: make reviewer, err: %v", err)

		return nil, fmt.Errorf("make reviewer, err: %v", err)
	}

	ls.log.Infof("the new reviewer was added with tx: %s", txHash)

	return participant, nil
}

func (ls *LibrarySrv) UpdateValidator(
	ctx context.Context,
	web3Address,
	emailAddress,
	name,
	middlename,
	surname,
	orcid,
	language string,
	sciences []string,
) (*storage.Participant, error) {
	// try to get the author record
	validatorResp, err := ls.GetValidator(ctx, web3Address)
	if err != nil {
		return nil, err
	}

	if emailAddress != "" {
		// it needs verification to be added here
		validatorResp.ValidatorInfo.EmailAddress = emailAddress
	}

	if name != "" {
		validatorResp.ValidatorInfo.Name = name
	}

	if middlename != "" {
		validatorResp.ValidatorInfo.MiddleName = middlename
	}

	if surname != "" {
		validatorResp.ValidatorInfo.Surname = surname
	}

	if orcid != "" {
		validatorResp.ValidatorInfo.Orcid = orcid
	}

	if language != "" {
		validatorResp.ValidatorInfo.Language = language
	}

	if len(sciences) > 0 {
		validatorResp.ValidatorInfo.Sciences = sciences
	}

	// update the Auhtor info in the storage(MongoDB)
	if err = ls.storage.UpdateValidatorInfo(ctx, validatorResp.ValidatorInfo); err != nil {
		ls.log.Errorf("UpdateValidator: error update validator info, err: %v", err)

		return nil, err
	}

	return validatorResp.BasicInfo, nil
}

// Publish work
// 1. save to the storage
// 2. publish to the IPFS
// 3. publish fingerPrint to the Library.sol
// returns: publish txId, error
func (ls *LibrarySrv) PublishWork(ctx context.Context, authorAddress string, work *storage.Work) (*storage.WorkResponse, string, error) {
	// check for the existence of the participant
	participant, err := ls.storage.GetParticipantByAddress(authorAddress)
	if err != nil {
		ls.log.Errorf("PublishWork: error get participant with address %s, err: %v", authorAddress, err)

		return nil, "", err
	}

	if participant.Role < storage.AuthorRole {
		return nil, "", fmt.Errorf("the participant nether author or validator")
	}

	// create work in Mongo and PostgreSQL databases
	workID, err := ls.storage.CreateWork(ctx, participant.ID, work)
	if err != nil {
		ls.log.Errorf("PublishWork: error create work, err: %v", err)

		return nil, "", fmt.Errorf("while creating a new work, err: %v", err)
	}

	fmt.Println("CREATED WORK ID: ", workID)
	// get the new work
	workResp, err := ls.storage.GetWorkByID(ctx, workID)
	if err != nil {
		ls.log.Errorf("PublishWork: error get work by id %s, err: %v", workID, err)

		return nil, "", err
	}

	txHash, err := ls.contractorSrv.PublishWork(ctx, &contractor.PublishWorkRequest{
		Authors: []string{authorAddress},
		Name:    work.Name,
		Uri:     "DUMMY URI",
		WorkId:  uuidToUint256(workResp.Work.ID),
		Price:   faucetCount,
	})
	if err != nil {
		ls.log.Errorf("PublishWork: publish work via contractor, err: %v", err)
	}

	return workResp, txHash.TxHash, nil
}

// PurchaseWork ...
func (ls *LibrarySrv) PurchaseWork(ctx context.Context, readerAddress, workID string, contract bool) error {
	// check for the existence of the participant
	participant, err := ls.storage.GetParticipantByAddress(readerAddress)
	if err != nil {
		ls.log.Errorf("PurchaseWork: error get participant with address %s, err: %v", readerAddress, err)

		return err
	}
	// get work by id
	work, err := ls.storage.GetWorkByID(ctx, workID)
	if err != nil {
		ls.log.Errorf("PurchaseWork: error get work by id %s, err: %v", workID, err)

		return fmt.Errorf("haven't got the work with id: %s", workID)
	}

	// check if he has already purchased the work
	if ls.storage.PurchasedWorkOrNot(participant.ID, workID) {
		return fmt.Errorf("you have already purchased this work")
	}

	if !contract {
		// burn some tokens from the buyer address
		// mint some token to the author address
		purchaseWorkResp, err := ls.contractorSrv.PurchaseWork(ctx, &contractor.PurchaseWorkRequest{
			WorkId:        workID,
			ReaderAddress: participant.Web3Address,
			AuthorAddress: work.Author.BasicInfo.Web3Address,
			Price:         work.Work.Price,
		})
		if err != nil {
			// todo errors.As()
			if !strings.Contains(err.Error(), "insufficient") {
				ls.log.Errorf("PurchaseWork: insufficient work, err: %v", err)
			}

			return err
		}

		ls.log.Infof("readerTxHash: %s| authorTxHash: %s", purchaseWorkResp.ReaderTxStatus.TxHash, purchaseWorkResp.AuthorTxStatus.TxHash)

	}
	// create work in Mongo and PostgreSQL databases
	if err = ls.storage.PurchaseWork(participant.ID, workID); err != nil {
		ls.log.Errorf("PurchaseWork: error purchase, participant id %s work id %s, err: %v", participant.ID, workID, err)

		return fmt.Errorf("while buying the work, err: %v", err)
	}

	return nil
}

// PurchaseWork ...
func (ls *LibrarySrv) PurchasedWorks(ctx context.Context, readerAddress string) ([]*storage.WorkResponse, error) {
	// check for the existence of the participant
	works, err := ls.storage.GetPurchasedWorks(ctx, readerAddress)
	if err != nil {
		ls.log.Errorf("PurchaseWork: error get pending works, err: %v", err)

		return nil, err
	}

	return works, nil
}

func (ls *LibrarySrv) GetPendingWorks(ctx context.Context) ([]*storage.WorkResponse, error) {
	// check for the existence of the participant
	works, err := ls.storage.GetPendingWorks(ctx)
	if err != nil {
		ls.log.Errorf("GetPendingWorks: error get pending works, err: %v", err)

		return nil, err
	}

	return works, nil
}

func (ls *LibrarySrv) GetAllWorks(ctx context.Context, readerAddress string) ([]*storage.WorkResponse, error) {
	// check for the existence of the participant
	works, err := ls.storage.GetAllWorks(ctx, readerAddress)
	if err != nil {
		ls.log.Errorf("GetAllWorks: error get all works, err: %v", err)

		return nil, err
	}

	return works, nil
}

func (ls *LibrarySrv) GetWorksByKeyWords(ctx context.Context, readerAddress string, keyWords []string) ([]*storage.WorkResponse, error) {
	// check for the existence of the participant
	works, err := ls.storage.GetWorkByKeyWords(ctx, readerAddress, keyWords)
	if err != nil {
		ls.log.Errorf("GetWorksByKeyWords: error get works by keywords %v, err: %v", keyWords, err)

		return nil, err
	}

	return works, nil
}

func (ls *LibrarySrv) GetWorkByID(ctx context.Context, authorAddress, workID string) (*storage.WorkResponse, error) {
	// check for the existence of the participant
	work, err := ls.storage.GetWorkByID(ctx, workID)
	if err != nil {
		ls.log.Errorf("GetWorkByID: error get work by id %s, err: %v", workID, err)

		return nil, err
	}

	return work, nil
}

func (ls *LibrarySrv) GetWorksByAuthorAddress(ctx context.Context, readerAddress, authorAddress string) ([]*storage.WorkResponse, error) {
	// check for the existence of the participant
	works, err := ls.storage.GetWorksByAuthorAddress(ctx, readerAddress, authorAddress)
	if err != nil {
		ls.log.Errorf("GetWorksByAuthorAddress: error get work by reader and author addresses %s, %s, err: %v", readerAddress, authorAddress, err)

		return nil, err
	}

	return works, nil
}

// Publish work
// 1. save to the storage
// 2. publish to the IPFS
// 3. publish fingerPrint to the Library.sol
// returns: publish txId, error
func (ls *LibrarySrv) ApproveWork(ctx context.Context, workID string) error {
	// check for the existence of the participant
	if err := ls.storage.ApproveWork(ctx, workID); err != nil {
		return fmt.Errorf("haven't got the participant with address %s", err)
	}

	// TODO update status in the Library.sol smart-contract

	// // publish to the IPFS
	// id := pinata.PublishJson(work)
	// fmt.Println("ID in IPFS: ", id)
	// publish fingerPrint to the Library.sol
	// author
	// id
	return nil

	// // publish to the IPFS
	// id := pinata.PublishJson(work)
	// fmt.Println("ID in IPFS: ", id)
	// // publish fingerPrint to the Library.sol
	// // author
	// // id
	// return id, "TX_ID", nil
}

func (ls *LibrarySrv) RemoveWork(ctx context.Context, workID string) error {
	// check for the existence of the participant
	if err := ls.storage.RemoveWork(ctx, workID); err != nil {
		ls.log.Errorf("RemoveWork: error remove work by work id %s, err: %v", workID, err)

		return err
	}
	// TODO update status in the Library.sol smart-contract
	return nil
}

func (ls *LibrarySrv) GetParticipantIdAndRoleByWeb3Address(address string) (string, storage.ParticipantRole) {
	participant, err := ls.storage.GetParticipantByAddress(address)
	if err != nil {
		ls.log.Errorf("GetParticipantIdAndRoleByWeb3Address: error get participant with address %s, err: %v", address, err)

		return "", -1
	}

	return participant.ID, participant.Role
}

func (ls *LibrarySrv) GetParticipantByWeb3Address(address string) (*storage.Participant, error) {
	participant, err := ls.storage.GetParticipantByAddress(address)
	if err != nil {
		ls.log.Errorf("GetParticipantIdAndRoleByWeb3Address: error get participant with address %s, err: %v", address, err)

		return nil, err
	}

	return participant, nil
}

func (ls *LibrarySrv) CreateBookmark(readerAddress, workID string) error {
	// get the participant by his address
	participant, err := ls.storage.GetParticipantByAddress(readerAddress)
	if err != nil {
		ls.log.Errorf("CreateBookmark: error get participant with address %s, err: %v", readerAddress, err)

		return err
	}
	// get the work by id
	work, err := ls.storage.GetParticipantWorkByID(workID)
	if err != nil {
		ls.log.Errorf("CreateBookmark: error get participant work with id %s, err: %v", workID, err)

		return err
	}

	return ls.storage.CreateBookmark(participant.ID, work.WorkID)
}

func (ls *LibrarySrv) GetBookmarksOf(ctx context.Context, readerAddress string) ([]*storage.WorkResponse, error) {
	// get the participant by his address
	participant, err := ls.storage.GetParticipantByAddress(readerAddress)
	if err != nil {
		ls.log.Errorf("GetBookmarksOf: error get participant with address %s, err: %v", readerAddress, err)

		return nil, err
	}
	// get the participant's bookmarks
	return ls.storage.GetBookmarksByParticipantID(ctx, participant.ID)
}

func (ls *LibrarySrv) RemoveBookmark(readerAddress, workID string) error {
	// get the participant by his address
	participant, err := ls.storage.GetParticipantByAddress(readerAddress)
	if err != nil {
		ls.log.Errorf("RemoveBookmark: error get participant with address %s, err: %v", readerAddress, err)

		return err
	}
	// get the work by id
	work, err := ls.storage.GetParticipantWorkByID(workID)
	if err != nil {
		ls.log.Errorf("RemoveBookmark: error get participant work with id %s, err: %v", workID, err)

		return err
	}

	return ls.storage.RemoveBookmark(participant.ID, work.WorkID)
}

/*//////////////////////////

//////// WORK Evaluation ////////

////////*/ ////////////////

func (ls *LibrarySrv) CreateOrUpdateWorkReview(ctx context.Context, validatorAddress string, review *storage.WorkReview) (*storage.WorkReview, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(validatorAddress)
	if err != nil {
		ls.log.Errorf("CreateOrUpdateWorkReview: error get participant with address %s, err: %v", validatorAddress, err)

		return nil, err
	}

	if participant.Role < storage.ValidatorRole {
		return nil, fmt.Errorf("the participant is not validator")
	}

	work, err := ls.storage.GetWorkByID(ctx, review.WorkID)
	if err != nil {
		ls.log.Errorf("CreateOrUpdateWorkReview: error get work by id, err: %v", err)

		return nil, fmt.Errorf("while get work by id, err: %v", err)
	}

	if work.Author.BasicInfo.ID == participant.ID {
		return nil, ErrValidationNotAllowed
	}

	review, updateRrr := ls.storage.UpdateOrCreateWorkReview(ctx, participant.ID, review)
	if updateRrr != nil {
		ls.log.Errorf("CreateOrUpdateWorkReview: error update or create work review, err: %v", err)

		return nil, fmt.Errorf("while creating validator, err: %v", err)
	}

	return review, nil
}

// GetWorkReviewByWorkID ...
func (ls *LibrarySrv) GetValidatorWorkReviewByWorkID(ctx context.Context, validatorAddress, workID string) (*storage.WorkReview, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(validatorAddress)
	if err != nil {
		ls.log.Errorf("GetValidatorWorkReviewByWorkID: error get participant with address %s, err: %v", validatorAddress, err)

		return nil, err
	}

	if participant.Role < storage.ValidatorRole {
		return nil, fmt.Errorf("the participant is not validator")
	}

	review, err := ls.storage.GetReviewByValidatorAndWorkID(ctx, participant.ID, workID)
	if err != nil {
		ls.log.Errorf("GetValidatorWorkReviewByWorkID: error get review, err: %v", err)

		return nil, fmt.Errorf("while getting the workReview, err: %v", err)
	}

	return review, nil
}

// GetWorkReviewsByWorkID ...
func (ls *LibrarySrv) GetWorkReviewsByWorkID(ctx context.Context, authorAddress, workID string) ([]*storage.WorkReview, error) {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(authorAddress)
	if err != nil {
		return nil, fmt.Errorf("haven't got the participant with address %s", authorAddress)
	}

	if participant.Role < storage.AuthorRole {
		return nil, fmt.Errorf("the participant is not author")
	}

	reviews, err := ls.storage.GetReviewByAuthorAndWorkID(ctx, participant.ID, workID)
	if err != nil {
		return nil, fmt.Errorf("while getting the work's reviews, err: %v", err)
	}

	return reviews, nil
}

func (ls *LibrarySrv) SubmitWorkReview(ctx context.Context, validatorAddress, workID string, status storage.WorkReviewStatus) error {
	// get the current participant and vefiry his role, status
	participant, err := ls.storage.GetParticipantByAddress(validatorAddress)
	if err != nil {
		ls.log.Errorf("SubmitWorkReview: error get participant with address %s, err: %v", validatorAddress, err)

		return err
	}

	if participant.Role < storage.ValidatorRole {
		return fmt.Errorf("the participant is not validator")
	}

	reviews, err := ls.storage.GetReviewsByWorkId(workID)
	if err != nil {
		ls.log.Errorf("SubmitWorkReview: error get all reviews with id %s, err: %v", workID, err)

		return err
	}

	if len(reviews) == 0 {
		return ErrNoReviews
	}

	var (
		participantsReview *storage.ParticipantsWorkReview
		isLastReview       bool
	)
	for _, r := range reviews {
		if r.Status == storage.WorkReviewSubmitted {
			isLastReview = true
		}

		if r.WorkID == workID && participant.ID == r.ParticipantID {
			participantsReview = r
		}
	}

	participantsReview.Status = status

	if err = ls.storage.SubmitWorkReview(ctx, participantsReview); err != nil {
		ls.log.Errorf("SubmitWorkReview: error submit work review, err: %v", err)

		return fmt.Errorf("while submitting the work review, err: %v", err)
	}

	if isLastReview {
		switch status {
		case storage.WorkReviewSubmitted:
			confirmErr := ls.storage.ConfirmWork(ctx, workID)
			if confirmErr != nil {
				ls.log.Errorf("SubmitWorkReview: error confirm work with id %s, err: %v", workID, err)

				return err
			}

		case storage.WorkReviewRejected, storage.WorkReviewSkipped:
			declinedErr := ls.storage.DeclineWork(ctx, workID)
			if declinedErr != nil {
				ls.log.Errorf("SubmitWorkReview: error confirm work with id %s, err: %v", workID, err)

				return err
			}
		}
	}

	return nil
}

// Faucet ...
func (ls *LibrarySrv) Faucet(account string) (string, error) {
	purchaseWorkResp, err := ls.contractorSrv.Faucet(context.Background(), &contractor.FaucetRequest{
		Address: account,
		Amount:  faucetCount,
	})
	if err != nil {
		return "", err
	}

	return purchaseWorkResp.TxHash, nil
}

func buildWorkResponse(inputWork storage.WorkResponse) (work storage.WorkResponse) {
	return storage.WorkResponse{}
}

func (ls *LibrarySrv) Stop() {}
