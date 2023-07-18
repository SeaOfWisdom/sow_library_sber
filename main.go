package main

import (
	"github.com/SeaOfWisdom/sow_library/src/common"
	"github.com/SeaOfWisdom/sow_library/src/config"
	"github.com/SeaOfWisdom/sow_library/src/container"
	"github.com/SeaOfWisdom/sow_library/src/rest-service"
	"github.com/SeaOfWisdom/sow_library/src/server"
	srv "github.com/SeaOfWisdom/sow_library/src/service"
)

// @title           SOW library API
// @version         1.0
// @description		Specification of interaction with the application

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @accept json
// @produce json
// @schemes http
func main() {
	di := container.CreateContainer()
	container.MustInvoke(di, func(
		config *config.Config,
		service *srv.LibrarySrv,
		restService *rest.RestSrv,
		grpcServer *server.GrpcServer,
	) {
		/* start services */
		grpcServer.Start()
		service.Start()
		restService.Start()

		/* wait for application termination */
		common.WaitForSignal()
		service.Start()
		restService.Stop()
		grpcServer.Stop()
	})
}

// cfg := config.NewConfig()

// strSrv := storage.NewStorageSrv(cfg, nil)
// work1 := storage.Work{
// 	ID:          "XXX",
// 	Name:        "1_XXX_NAME",
// 	Description: "1_XXX_DESCO",
// 	Annotation:  "1_XXX_CCC",
// 	AuthorID:    "1",
// 	Tags:        []string{"science, sport, facilities"},
// }
// fmt.Println(strSrv.PutWork(&work1))

// works, err := strSrv.GetWorkByAuthor("1")
// if err != nil {
// 	panic(err)
// }

// for _, work := range works {
// 	fmt.Println(work)
// }

// return

// pinata.PublishJson()
// return
// // create service
// instance := w3scli.CreateW3sSrv()
// put example json work
// id, err := instance.PutJsonFile()
// if err != nil {
// 	fmt.Println("while putting the json file")
// 	panic(err)
// }
// fmt.Println("ID: ", id)
// id := "bafybeibrbbsk237jfordm5gezfyr2ejn54wfkqvgncddvkahvdwzktno6m"

// // try to get it back and read
// _, err := instance.GetStatusForCid(id)
// if err != nil {
// 	fmt.Println("while getting a status for the cid")
// 	panic(err)
// }

// fmt.Println(instance.GetFiles(id))
