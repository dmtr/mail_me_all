package rpc

import (
	fmt "fmt"

	"github.com/dmtr/mail_me_all/backend/config"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//GetRpcConection - returns connection to grpc server
func GetRpcConection(conf *config.Config) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if conf.PemFile != "" {
		creds, err := credentials.NewClientTLSFromFile(conf.PemFile, "")
		if err != nil {
			log.Fatalf("Can read credentials file %s:, got error %s", conf.PemFile, err)
		}
		log.Info("Load credentials")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	serverAddr := fmt.Sprintf("%s:%d", conf.TwProxyHost, conf.TwProxyPort)
	return grpc.Dial(serverAddr, opts...)
}

//GetRpcClient - returns grpc client
func GetRpcClient(conn *grpc.ClientConn) TwProxyServiceClient {
	return NewTwProxyServiceClient(conn)
}
