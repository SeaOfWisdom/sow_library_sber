package srv

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/SeaOfWisdom/sow_library/src/config"
	contractor "github.com/SeaOfWisdom/sow_proto/contractor-srv"
)

func (ls *LibrarySrv) MigrateFromMongo() {
	ctx := context.Background()
	// 1. get all participants
	participants := ls.storage.GetAllParticipants()

	// 2. get all works from admin address
	papers, err := ls.storage.GetAllWorks(ctx, config.AdminAddresses["chillhacker"])
	if err != nil {
		panic(fmt.Sprintf("MigrateFromMongo: error while getting all works, err: %v", err))
	}

	// 3. ensure all participants are added into smart contract
	for _, participant := range participants {
		fmt.Println("participant.Web3Address: ", participant.Web3Address)
		role, err := ls.contractorSrv.GetParticipantRole(ctx, &contractor.AccountRequest{
			Address: participant.Web3Address,
		})
		if err != nil {
			ls.log.Errorf("MigrateFromMongo: get participant role, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if role.Role != 0 {
			continue
		}

		txHash, err := ls.contractorSrv.AddParticipant(ctx, &contractor.AccountRequest{
			Address: participant.Web3Address,
		})
		if err != nil {
			ls.log.Errorf("MigrateFromMongo: add participant, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		ls.log.Infof("participant(%s) with address(%s) was added with tx %s", participant.NickName, participant.Web3Address, txHash)
		time.Sleep(5 * time.Second)
	}

	// 4. ensure all works are added into smart contract
	for _, paper := range papers {
		paperAddress, err := ls.contractorSrv.GetPaperById(ctx, &contractor.PaperByIdRequest{
			Id: uuidToUint256(paper.Work.ID),
		})
		if err != nil {
			ls.log.Errorf("MigrateFromMongo: get paper by id, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if paperAddress.Address != "0x0000000000000000000000000000000000000000" {
			continue
		}

		// check the author current role
		// if he is not author -> make him author
		role, err := ls.contractorSrv.GetParticipantRole(ctx, &contractor.AccountRequest{
			Address: paper.Author.BasicInfo.Web3Address,
		})
		if err != nil {
			ls.log.Errorf("MigrateFromMongo: get participant role, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if role.Role < 2 {
			txHash, err := ls.contractorSrv.MakeAuthor(ctx, &contractor.AccountRequest{
				Address: paper.Author.BasicInfo.Web3Address,
			})
			if err != nil {
				ls.log.Errorf("MigrateFromMongo: make author, err: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			ls.log.Infof("participants(%s) became author with tx %s", paper.Author.BasicInfo.Web3Address, txHash)
			time.Sleep(5 * time.Second)
		}

		txHash, err := ls.contractorSrv.PublishWork(ctx, &contractor.PublishWorkRequest{
			Authors: []string{paper.Author.BasicInfo.Web3Address},
			Name:    paper.Work.Name,
			Uri:     "DUMMY_URI",
			WorkId:  uuidToUint256(paper.Work.ID),
			Price:   faucetCount,
		})
		if err != nil {
			ls.log.Errorf("MigrateFromMongo: add participant, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		ls.log.Infof("paper(%s) was added with tx %s", paper.Work.Name, txHash)
		time.Sleep(5 * time.Second)
	}
}

func uuidToUint256(uuid string) string {
	var i big.Int
	i.SetString(strings.Replace(uuid, "-", "", 4), 16)
	return i.String()
}
