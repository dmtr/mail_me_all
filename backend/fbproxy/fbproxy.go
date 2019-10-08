package fbproxy

import (
	"context"

	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	log "github.com/sirupsen/logrus"
)

//ServiceServer - grpc service
type ServiceServer struct {
	facebook fbwrapper.Facebook
}

//NewServiceServer - returns new service server instance
func NewServiceServer(f fbwrapper.Facebook) *ServiceServer {
	return &ServiceServer{facebook: f}
}

//GetAccessToken - returns long lived access token
func (s *ServiceServer) GetAccessToken(ctx context.Context, user *pb.UserToken) (*pb.UserToken, error) {
	userID, err := s.facebook.VerifyFbToken(user.AccessToken)
	if err != nil {
		log.Errorf("Invalid access token: error %s", err)
		return nil, err
	}

	res, err := s.facebook.GenerateLongLivedToken(user.AccessToken)
	if err != nil {
		log.Errorf("Can not create long lived token, error: %s", err)
		return nil, err
	}

	u := pb.UserToken{
		UserId:      userID,
		AccessToken: res.AccessToken,
		ExpiresIn:   uint64(res.ExpiresIn),
	}

	return &u, nil
}

//GetUserInfo - returns fb user info
func (s *ServiceServer) GetUserInfo(ctx context.Context, user *pb.UserToken) (*pb.UserInfo, error) {
	u, err := s.facebook.GetUserInfo(user.UserId, user.AccessToken)
	if err != nil {
		log.Errorf("Can not get user %s info, error: %s", user.UserId, err)
		return nil, err
	}

	res := pb.UserInfo{
		UserId: u.UserID,
		Name:   u.FirstName,
		Email:  u.Email,
	}

	return &res, nil
}
