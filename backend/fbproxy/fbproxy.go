package fbproxy

import (
	"context"

	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	log "github.com/sirupsen/logrus"
)

//ServiceServer - grpc service
type ServiceServer struct {
	app *app.App
}

//NewServiceServer - returns new service server instance
func NewServiceServer(app *app.App) *ServiceServer {
	return &ServiceServer{app: app}
}

//GetAccessToken - returns long lived access token
func (s *ServiceServer) GetAccessToken(ctx context.Context, user *pb.NewUser) (*pb.User, error) {
	userID, err := fbwrapper.VerifyFbToken(user.AccessToken, s.app.Conf.FbAppID, s.app.Conf.AppSecret, s.app.Conf.FbRedirectURI)
	if err != nil {
		log.Errorf("Invalid access token: error %s", err)
		return nil, err
	}

	res, err := fbwrapper.GenerateLongLivedToken(user.AccessToken, s.app.Conf.FbAppID, s.app.Conf.AppSecret)
	if err != nil {
		log.Errorf("Can not create long lived token, error: %s", err)
		return nil, err
	}

	u := pb.User{
		UserId:      userID,
		AccessToken: res.AccessToken,
		ExpiresIn:   uint64(res.ExpiresIn),
	}

	return &u, nil
}