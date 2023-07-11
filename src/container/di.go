package container

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SeaOfWisdom/sow_library/src/config"
	"github.com/SeaOfWisdom/sow_library/src/log"
	"github.com/SeaOfWisdom/sow_library/src/server"
	lib "github.com/SeaOfWisdom/sow_library/src/service"
	contractorProto "github.com/SeaOfWisdom/sow_proto/contractor-srv"
	jwtProto "github.com/SeaOfWisdom/sow_proto/jwt-srv"
	ocrProto "github.com/SeaOfWisdom/sow_proto/ocr-srv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SeaOfWisdom/sow_library/src/rest-service"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"go.uber.org/dig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateContainer() *dig.Container {
	container := dig.New()
	/* init base */
	must(container.Provide(config.NewConfig))
	must(container.Provide(log.NewLogger))
	must(container.Provide(func() *emitter.Emitter {
		return emitter.New(10)
	}))
	/* init mongo */
	must(container.Provide(func(config *config.Config) *mongo.Database {
		client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoAddress).SetAuth(options.Credential{
			Username: config.MongoUser,
			Password: config.MongoPassword,
		}))
		if err != nil {
			panic(fmt.Errorf("while creating a client for MongoDB, err: %v", err))
		}
		if err := client.Connect(context.Background()); err != nil {
			panic(fmt.Errorf("while connecting to MongoDB, err: %v", err))
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Ping(ctx, readpref.Nearest()); err != nil {
			panic(fmt.Errorf("while connecting to MongoDB(ping to client), err: %v", err))
		}
		return client.Database(config.MongoDBName)
	}))
	/* init postgres */
	must(container.Provide(func(config *config.Config) *gorm.DB {
		opts := strings.Split(config.PostgresAddress, ":")
		if len(opts) != 2 {
			panic(fmt.Errorf("wrong PostgresAddress: %s", config.PostgresAddress))
		}
		postgresOpts := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			opts[0], opts[1], config.PostgresUser, config.PostgresDbName, config.PostgresPassword)
		gormCli, err := gorm.Open(postgres.Open(postgresOpts), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic(fmt.Errorf("failed to set connection to PostgreSQL(%s), err: %v ", postgresOpts, err.Error()))
		}
		return gormCli
	}))
	/* set connection to the internal services */
	/* JWT Service */
	must(container.Provide(func(config *config.Config) jwtProto.JwtServiceClient {
		conn, err := grpc.Dial(config.JWTServiceGRpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("unable to start JWT Service GRPC client: %v", err))
		}
		return jwtProto.NewJwtServiceClient(conn)
	}))
	/* OCR Service */
	must(container.Provide(func(config *config.Config) ocrProto.OCRClient {
		conn, err := grpc.Dial(config.OCRServiceGRpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("unable to start OCR Service GRPC client: %v", err))
		}
		return ocrProto.NewOCRClient(conn)
	}))
	/* Contractor Service */
	must(container.Provide(func(config *config.Config) contractorProto.ContractorServiceClient {
		conn, err := grpc.Dial(config.ContractorServiceGRpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("unable to start Contractor Service GRPC client: %v", err))
		}
		return contractorProto.NewContractorServiceClient(conn)
	}))
	/* initialize the storage consists of mongoDB and postgreDB */
	must(container.Provide(storage.NewStorageSrv))
	/* initialize internal services */
	must(container.Provide(lib.NewLibrarySrv))
	must(container.Provide(rest.NewRestSrv))
	must(container.Provide(server.NewGrpcServer))
	return container
}

func MustInvoke(container *dig.Container, function interface{}, opts ...dig.InvokeOption) {
	must(container.Invoke(function, opts...))
}

func must(err error) {
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DI: %s", err))
	}
}
