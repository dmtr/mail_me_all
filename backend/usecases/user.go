package usecases

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/errors"
	"github.com/dmtr/mail_me_all/backend/mail"
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const ConfirmationEmailSubj = "Email address confirmation"

type JWTClaims struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// UserUseCase implementation
type UserUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.TwProxyServiceClient
	Conf          *config.Config
	Tmpl          *template.Template
}

// NewUserUseCase implementation
func NewUserUseCase(datastore models.UserDatastore, client pb.TwProxyServiceClient, conf *config.Config) *UserUseCase {
	tmpl := template.Must(template.New("confirm.html").ParseFiles(filepath.Join(conf.TemplatePath, "confirm.html")))
	return &UserUseCase{UserDatastore: datastore, RpcClient: client, Conf: conf, Tmpl: tmpl}
}

// GetUserByID implementation
func (u UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.UserDatastore.GetUser(ctx, userID)
	if err != nil {
		return user, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}
	return user, err
}

// SignInWithTwitter implementation
func (u UserUseCase) SignInWithTwitter(ctx context.Context, twitterID, name, email, screenName, accessToken, tokenSecret string) (models.User, error) {
	user := models.User{Name: name}
	var profileUrl string

	req := pb.UserInfoRequest{
		TwitterId:    twitterID,
		AccessToken:  accessToken,
		AccessSecret: tokenSecret,
		ScreenName:   screenName,
	}
	if userInfo, err := u.RpcClient.GetUserInfo(context.Background(), &req); err != nil {
		log.Errorf("Can not get user info: %s", err)
	} else {
		user.Email = userInfo.Email
		profileUrl = userInfo.ProfileImageUrl
	}

	var userExists bool

	twitterUser, err := u.UserDatastore.GetTwitterUserByID(ctx, twitterID)
	if err != nil {
		code := errors.GetErrorCode(err)
		if code == errors.NotFound {
			userExists = false
		} else {
			return user, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		userExists = true
	}

	if userExists {
		user.ID = twitterUser.UserID

		if user, err = u.UserDatastore.UpdateUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}

		twitterUser.AccessToken = accessToken
		twitterUser.TokenSecret = tokenSecret
		twitterUser.ProfileIMGURL = profileUrl

		if _, err = u.UserDatastore.UpdateTwitterUser(ctx, twitterUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
		log.Debugf("New user %s", user)

		twUser := models.TwitterUser{
			UserID:        user.ID,
			TwitterID:     twitterID,
			AccessToken:   accessToken,
			TokenSecret:   tokenSecret,
			ProfileIMGURL: profileUrl,
		}

		if _, err = u.UserDatastore.InsertTwitterUser(ctx, twUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}

	}
	return user, err

}

func (u UserUseCase) SearchTwitterUsers(ctx context.Context, userID uuid.UUID, query string) ([]models.TwitterUserSearchResult, error) {
	twitterUser, err := u.UserDatastore.GetTwitterUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	req := pb.UserSearchRequest{
		TwitterId:    twitterUser.TwitterID,
		AccessToken:  twitterUser.AccessToken,
		AccessSecret: twitterUser.TokenSecret,
		Query:        query,
	}

	res, err := u.RpcClient.SearchUsers(context.Background(), &req)
	if err != nil {
		log.Errorf("Can not find users: %s", err)
		return nil, err
	}

	users := make([]models.TwitterUserSearchResult, 0, len(res.Users))
	for _, user := range res.Users {
		u := models.TwitterUserSearchResult{
			TwitterID:     user.TwitterId,
			Name:          user.Name,
			ScreenName:    user.ScreenName,
			ProfileIMGURL: user.ProfileImageUrl,
		}
		users = append(users, u)
	}

	return users, err
}

func (u UserUseCase) AddSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	_, err := u.UserDatastore.GetUser(ctx, subscription.UserID)
	if err != nil {
		return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	s, err := u.UserDatastore.InsertSubscription(ctx, subscription)
	if err != nil {
		return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	userEmail := models.UserEmail{
		UserID: subscription.UserID,
		Email:  subscription.Email,
		Status: models.EmailStatusNew,
	}

	email, err := u.UserDatastore.InsertUserEmail(ctx, userEmail)
	sendConfirmationEmail := false
	if err != nil {
		e := err.(*db.DbError)
		if e.IsUniqueViolationError() {
			err = nil
			log.Infof("Record with email %s exists %s", userEmail, email)
		}
	} else {
		sendConfirmationEmail = true
	}

	if err != nil {
		return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	if sendConfirmationEmail {
		if err = u.sendConfirmationEmail(userEmail.Email, userEmail.UserID.String()); err != nil {
			return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	}

	return s, nil
}

func (u UserUseCase) GetToken(email, userID string) (string, error) {
	exp := time.Now().Unix() + 24*60*60
	claims := JWTClaims{
		email,
		userID,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(u.Conf.EncryptKey))
	if err != nil {
		return "", err
	}

	return ss, err
}

func (u UserUseCase) getEmailConfirmationLink(email, userID string) (string, error) {
	token, err := u.GetToken(email, userID)
	if err != nil {
		return "", err
	}

	link := &url.URL{
		Scheme:   "https",
		Host:     u.Conf.Domain,
		Path:     "confirm-email",
		RawQuery: fmt.Sprintf("token=%s", token),
	}

	return link.String(), err
}

func (u UserUseCase) sendConfirmationEmail(email, userID string) error {
	var buf strings.Builder
	link, err := u.getEmailConfirmationLink(email, userID)
	if err != nil {
		return err
	}

	err = u.Tmpl.Execute(&buf, link)
	if err != nil {
		return err
	}

	err = mail.SendEmail(u.Conf.MgDomain, u.Conf.MgAPIKEY, u.Conf.From, email, ConfirmationEmailSubj, buf.String())
	return err
}

func (u UserUseCase) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	s, err := u.UserDatastore.GetSubscriptions(ctx, userID)
	if err != nil {
		return s, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return s, nil
}

func (u UserUseCase) UpdateSubscription(ctx context.Context, userID uuid.UUID, subscription models.Subscription) (models.Subscription, error) {
	if userID != subscription.UserID {
		err := fmt.Errorf("User %s can not edit subscription %s", userID, subscription)
		return subscription, NewUseCaseError(err.Error(), errors.AuthRequired)
	}

	userEmail := models.UserEmail{
		UserID: subscription.UserID,
		Email:  subscription.Email,
	}

	sendConfirmationEmail := false
	email, err := u.UserDatastore.GetUserEmail(ctx, userEmail)
	if err != nil {
		e := err.(*db.DbError)
		if e.HasNoRows() {
			sendConfirmationEmail = true
		} else {
			return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	}

	if email.UserID != subscription.UserID {
		log.Warningf("Email %s belongs to another user %s", subscription.Email, email)
		return subscription, NewUseCaseError("Email belongs to another user", errors.AuthRequired)
	}

	if sendConfirmationEmail || email.Status == models.EmailStatusNew {
		if err = u.sendConfirmationEmail(userEmail.Email, userEmail.UserID.String()); err != nil {
			return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	}

	s, err := u.UserDatastore.UpdateSubscription(ctx, subscription)
	if err != nil {
		return s, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return s, nil
}

func (u UserUseCase) DeleteSubscription(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error {
	subscription, err := u.UserDatastore.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	if userID != subscription.UserID {
		err := fmt.Errorf("User %s can not edit subscription %s", userID, subscription)
		return NewUseCaseError(err.Error(), errors.AuthRequired)
	}

	err = u.UserDatastore.DeleteSubscription(ctx, subscription)
	if err != nil {
		return NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return nil
}

func (u UserUseCase) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	err := u.UserDatastore.RemoveUser(ctx, userID)
	if err != nil {
		return NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return err
}

func (u UserUseCase) parseToken(token string) (models.UserEmail, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.Conf.EncryptKey), nil
	})

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		email := claims["email"]
		userID := claims["user_id"]
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			return models.UserEmail{}, err
		}
		return models.UserEmail{Email: email.(string), UserID: uid}, nil
	}
	return models.UserEmail{}, err
}

func (u UserUseCase) ConfirmEmail(ctx context.Context, token string) error {
	emailToConfirm, err := u.parseToken(token)
	if err != nil {
		return err
	}

	emailFromDB, err := u.UserDatastore.GetUserEmail(ctx, emailToConfirm)
	if err != nil {
		return err
	}

	if emailToConfirm.UserID != emailFromDB.UserID {
		log.Errorf("Emails not matching %s %s", emailToConfirm, emailFromDB)
		return NewUseCaseError("Can't confirm email", errors.AuthRequired)
	}

	emailToConfirm.Status = models.EmailStatusConfirmed
	_, err = u.UserDatastore.UpdateUserEmail(ctx, emailToConfirm)
	if err != nil {
		return err
	}

	return nil
}
