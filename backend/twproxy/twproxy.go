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
		TwitterId:       res.TwitterID,
		Name:            res.Name,
		Email:           res.Email,
		ScreenName:      res.ScreenName,
		ProfileImageUrl: res.ProfileIMGURL,
	}

	return &u, nil
}

//SearchUsers searches twitter users
func (s *ServiceServer) SearchUsers(ctx context.Context, request *pb.UserSearchRequest) (*pb.UserSearchResult, error) {
	users, err := s.twitter.SearchUsers(request.AccessToken, request.AccessSecret, request.TwitterId, request.Query)
	if err != nil {
		return nil, err
	}

	res := pb.UserSearchResult{
		Users: make([]*pb.UserInfo, 0, len(users)),
	}

	for _, user := range users {
		u := pb.UserInfo{
			TwitterId:       user.TwitterID,
			Name:            user.Name,
			Email:           user.Email,
			ScreenName:      user.ScreenName,
			ProfileImageUrl: user.ProfileIMGURL,
		}
		res.Users = append(res.Users, &u)
	}

	return &res, err
}

//
func (s *ServiceServer) GetUserTimeline(ctx context.Context, request *pb.UserTimelineRequest) (*pb.UserTimelineResponse, error) {
	tweets, err := s.twitter.GetUserTimeline(
		request.AccessToken, request.AccessSecret, request.TwitterId, request.ScreenName, request.SinceId, request.Count, request.IgnoreRt, request.IgnoreReplies)
	if err != nil {
		return nil, err
	}

	res := pb.UserTimelineResponse{
		Tweets: make([]*pb.Tweet, 0, len(tweets)),
	}

	for _, tweet := range tweets {
		t := pb.Tweet{
			IdStr:                tweet.IDStr,
			Text:                 tweet.Text,
			FullText:             tweet.FullText,
			InReplyToStatusIdStr: tweet.InReplyToStatusIDStr,
			InReplyToUserIdStr:   tweet.InReplyToUserIDStr,
			UserId:               tweet.UserID,
			UserName:             tweet.UserName,
			UserScreenName:       tweet.UserScreenName,
			UserProfileImageUrl:  tweet.UserProfileImageUrl,
		}
		res.Tweets = append(res.Tweets, &t)
	}

	return &res, err
}
