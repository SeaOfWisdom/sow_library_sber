package config

import (
	"github.com/namsral/flag"
)

var AdminAddresses = map[string]string{
	"chillhacker": "0xdd868980ef73edcbc1ff758f6e53023be18e2a52",
	"esaurov":     "0xb41649C7f675D77dbEf7D1018653f5f4eA9EeEE7",
}

type Config struct {
	/* gRPC */
	GrpcAddress string
	/* REST */
	RestAddress string
	/* Prometheus */
	PrometheusAddress string
	/* MongoDB */
	MongoAddress  string
	MongoUser     string
	MongoPassword string
	MongoDBName   string
	/* Postgres */
	PostgresAddress  string
	PostgresPort     string
	PostgresDriver   string
	PostgresDbName   string
	PostgresUser     string
	PostgresPassword string
	/* Cron */
	AddRewardsCron   string
	UpdateRewadsCron string
	/* Pinata */
	PinataURL string
	PinataJWT string
	/* Metric */
	MetricService     string
	MetricServiceGrpc string
	PartnersPercent   int64
	/* Internal communication services */
	JWTServiceGRpcAddress        string
	OCRServiceGRpcAddress        string
	ContractorServiceGRpcAddress string
}

func NewConfig() *Config {
	config := &Config{}
	/* gRPC */
	flag.StringVar(&config.GrpcAddress, "grpc-address", "0.0.0.0:8060", "gRPC address and port for inter-service communications")
	/* REST */
	flag.StringVar(&config.RestAddress, "rest-address", "0.0.0.0:8005", "REST address and port for public communications")
	/* Prometheus */
	flag.StringVar(&config.PrometheusAddress, "prometheus-address", "localhost:8075", "host and port for prometheus")
	/* Mongo */
	flag.StringVar(&config.MongoAddress, "mongo-address", "mongodb://localhost:27017", "")
	flag.StringVar(&config.MongoUser, "mongo-user", "root", "")
	flag.StringVar(&config.MongoPassword, "momgo-password-name", "simsim", "")
	flag.StringVar(&config.MongoDBName, "momgo-db-name", "pp", "")
	/* Postgres */
	flag.StringVar(&config.PostgresAddress, "postgres-address", "0.0.0.0:5432", "")
	flag.StringVar(&config.PostgresDbName, "postgres-db-name", "sow", "")
	flag.StringVar(&config.PostgresUser, "postgres-user", "postgres", "")
	flag.StringVar(&config.PostgresPassword, "postgres-password", "simsim", "")
	flag.StringVar(&config.PostgresDriver, "postgres-driver", "postgres", "")
	/* Cron */
	flag.StringVar(&config.AddRewardsCron, "add-rewards-cron", "*/1 * * * *", "")
	flag.StringVar(&config.UpdateRewadsCron, "update-rewards-cron", "*/3 * * * *", "")
	/* Pinata */
	flag.StringVar(&config.PinataURL, "pinata-url", "*/3 * * * *", "")
	flag.StringVar(&config.PinataJWT, "pinata-jwt", "*/3 * * * *", "")
	/* Internal communication services */
	flag.StringVar(&config.JWTServiceGRpcAddress, "jwt-service-address", "0.0.0.0:5304", "")
	flag.StringVar(&config.OCRServiceGRpcAddress, "ocr-service-address", "0.0.0.0:50051", "")
	flag.StringVar(&config.ContractorServiceGRpcAddress, "contractor-service-address", "0.0.0.0:5305", "")

	/* parse config from envs or config files */
	flag.Parse()
	return config
}
