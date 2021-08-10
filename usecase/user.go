package usecase

import (
	"net/http"

	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"
	"go-ddd/constant"
	"go-ddd/domain/entity"
	"go-ddd/domain/repository"
	"go-ddd/resource/request"
	"go-ddd/resource/response"
	"go-ddd/util"
	"go-ddd/util/xerrors"
	"gorm.io/gorm"
)

type IUser interface {
	Create(req *request.UserCreate) (uint, error)

	ResetPasswordRequest(req *request.UserResetPasswordRequest) (*response.UserResetPasswordRequest, error)
	ResetPassword(req *request.UserResetPassword) error
	Login(req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(refreshToken string) (*response.UserLogin, error)
}

type user struct {
	userRepo repository.IUser
}

func NewUser(tr repository.IUser) IUser {
	return &user{
		userRepo: tr,
	}
}

func (u user) Create(req *request.UserCreate) (uint, error) {
	verr := xerrors.NewValidation()

	email, err := u.userRepo.EmailExists(util.DB, req.Email)
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

	id, err := u.userRepo.Create(util.DB, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) ResetPasswordRequest(req *request.UserResetPasswordRequest) (*response.UserResetPasswordRequest, error) {
	user, err := u.userRepo.GetByEmail(util.DB, req.Email)
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

	err = util.DB.Transaction(
		func(tx *gorm.DB) error {
			err = u.userRepo.Update(util.DB, user)
			if err != nil {
				return err
			}

			err = sendMail(user.Email, "パスワードリセット", token)
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

func (u user) ResetPassword(req *request.UserResetPassword) error {
	verr := xerrors.NewValidation()

	user, err := u.userRepo.GetByRecoveryToken(util.DB, req.RecoveryToken)
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

	return u.userRepo.Update(util.DB, user)
}

func (u user) Login(req *request.UserLogin) (*response.UserLogin, error) {
	user, err := u.userRepo.GetByEmail(util.DB, req.Email)
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
