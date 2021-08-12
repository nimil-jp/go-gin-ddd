package usecase

import (
	"context"
	"net/http"

	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"
	"go-ddd/constant"
	"go-ddd/domain/entity"
	"go-ddd/domain/repository"
	"go-ddd/pkg/rdb"
	"go-ddd/pkg/xerrors"
	"go-ddd/resource/request"
	"go-ddd/resource/response"
)

type IUser interface {
	Create(ctx context.Context, req *request.UserCreate) (uint, error)

	ResetPasswordRequest(
		ctx context.Context,
		req *request.UserResetPasswordRequest,
	) (*response.UserResetPasswordRequest, error)
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(refreshToken string) (*response.UserLogin, error)
}

type user struct {
	emailRepo repository.IEmail
	userRepo  repository.IUser
}

func NewUser(email repository.IEmail, tr repository.IUser) IUser {
	return &user{
		emailRepo: email,
		userRepo:  tr,
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) (uint, error) {
	verr := xerrors.NewValidation()

	email, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return 0, err
	}

	if email {
		verr.Add("Email", "既に使用されています")
	}

	newUser, err := entity.NewUser(verr, req)
	if err != nil {
		return 0, err
	}

	if verr.IsInValid() {
		return 0, verr
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) ResetPasswordRequest(
	ctx context.Context,
	req *request.UserResetPasswordRequest,
) (*response.UserResetPasswordRequest, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		switch v := err.(type) {
		case *xerrors.Expected:
			if !v.ChangeStatus(http.StatusNotFound, http.StatusOK) {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	var token string
	var res response.UserResetPasswordRequest

	token, res.Duration, res.Expire, err = user.ResetPasswordRequest()
	if err != nil {
		return nil, err
	}

	err = rdb.Transaction(
		ctx,
		func(ctx context.Context) error {
			err = u.userRepo.Update(ctx, user)
			if err != nil {
				return err
			}

			err = u.emailRepo.Send(user.Email, "パスワードリセット", token)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u user) ResetPassword(ctx context.Context, req *request.UserResetPassword) error {
	verr := xerrors.NewValidation()

	user, err := u.userRepo.GetByRecoveryToken(ctx, req.RecoveryToken)
	if err != nil {
		return err
	}

	err = user.ResetPassword(verr, req)
	if err != nil {
		return err
	}

	if verr.IsInValid() {
		return verr
	}

	return u.userRepo.Update(ctx, user)
}

func (u user) Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.PasswordIsValid(req.Password) {
		var res response.UserLogin

		res.Token, res.RefreshToken, err = jwt.IssueToken(constant.DefaultRealm, jwt.Claims{})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return &res, nil
	}
	return nil, nil
}

func (u user) RefreshToken(refreshToken string) (*response.UserLogin, error) {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(constant.DefaultRealm, refreshToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !ok {
		return nil, nil
	}
	return &res, nil
}
