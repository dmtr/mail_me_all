package rpc

import (
	fmt "fmt"

	"github.com/dmtr/mail_me_all/backend/config"
	grpc "google.golang.org/grpc"
)

//GetRpcConection - returns connection to grpc server
func GetRpcConection(conf *config.Config) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	serverAddr := fmt.Sprintf("%s:%d", conf.FBProxyHost, conf.FBProxyPort)
	return grpc.Dial(serverAddr, opts...)
}

//GetRpcClient - returns grpc client
func GetRpcClient(conn *grpc.ClientConn) FbProxyServiceClient {
	return NewFbProxyServiceClient(conn)
}
