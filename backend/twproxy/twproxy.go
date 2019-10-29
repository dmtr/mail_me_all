package twproxy

import (
	"context"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/dmtr/mail_me_all/backend/twapi"
)

//ServiceServer - grpc service
type ServiceServer struct {
	twitter twapi.Twitter
}

//NewServiceServer - returns new service server instance
func NewServiceServer(t twapi.Twitter) *ServiceServer {
	return &ServiceServer{twitter: t}
}

//GetUserInfo - returns twitter user info
func (s *ServiceServer) GetUserInfo(ctx context.Context, request *pb.UserInfoRequest) (*pb.UserInfo, error) {
	res, err := s.twitter.GetUserInfo(request.AccessToken, request.AccessSecret, request.TwitterId, request.ScreenName)
	if err != nil {
		return nil, err
	}

	u := pb.UserInfo{
		TwitterId:  res.TwitterID,
		Name:       res.Name,
		Email:      res.Email,
		ScreenName: res.ScreenName,
	}

	return &u, nil
}
