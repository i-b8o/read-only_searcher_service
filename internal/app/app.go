package app

import (
	"context"
	"fmt"
	"net"
	"read-only_search/internal/config"
	"read-only_search/internal/controller/v1"
	postgressql "read-only_search/internal/data_providers/db/postgresql"
	"read-only_search/internal/domen/service"
	"read-only_search/pkg/client/postgresql"

	"time"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"

	"github.com/i-b8o/logging"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	logger     logging.Logger
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logger := logging.GetLogger(config.AppConfig.LogLevel)

	logger.Print("Postgres initializing")
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}
	docAdapter := postgressql.NewDocStorage(pgClient)
	chapterAdapter := postgressql.NewChapterStorage(pgClient)
	paragraphAdapter := postgressql.NewParagraphStorage(pgClient)
	generalAdapter := postgressql.NewParagraphStorage(pgClient)

	docService := service.NewDocService(docAdapter)
	chapterService := service.NewChapterService(chapterAdapter)
	paragraphService := service.NewParagraphsService(paragraphAdapter)
	generalService := service.NewGeneralService(generalAdapter)

	controller := controller.NewDocGRPCService(docService, chapterService, paragraphService, generalService)
	// read ca's cert, verify to client's certificate
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caPem, err := ioutil.ReadFile(homeDir + "/certs/ca-cert.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create cert pool and append ca's cert
	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(caPem) {
	// 	log.Fatal(err)
	// }

	// // read server cert & key
	// serverCert, err := tls.LoadX509KeyPair(homeDir+"/certs/server-cert.pem", homeDir+"/certs/server-key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // configuration of the certificate what we want to
	// conf := &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	ClientAuth:   tls.RequireAndVerifyClientCert,
	// 	ClientCAs:    certPool,
	// }

	// //create tls certificate
	// tlsCredentials := credentials.NewTLS(conf)

	// grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	// pb.RegisterReadOnlyDocGRPCServer(grpcServer, docGrpcService)
	logger.Print("grpc server initializing")
	grpcServer := grpc.NewServer()
	pb.RegisterSearcherGRPCServer(grpcServer, controller)

	return App{cfg: config, grpcServer: grpcServer, logger: logger}, nil
}

func (a *App) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	a.logger.Printf("started server on %s", address)
	return a.grpcServer.Serve(listener)
}
