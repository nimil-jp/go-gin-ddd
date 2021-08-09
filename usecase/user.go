package usecase

import (
	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"
	"go-ddd/constant"
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"go-ddd/resource/request"
	"go-ddd/resource/response"
	"go-ddd/util"
	"go-ddd/util/validate"
)

type IUser interface {
	Create(req *request.UserCreate) (uint, error)
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
	verr := validate.NewValidationError()

	email, err := u.userRepo.EmailExists(util.DB, req.Email)
	if err != nil {
		return 0, err
	}

	if email {
		verr.Add("email", "既に使われています")
	}

	if !verr.Validate(req) {
		return 0, verr
	}

	hashed, err := genHashedPassword(req.Password)
	if err != nil {
		return 0, err
	}

	id, err := u.userRepo.Create(
		util.DB, &model.User{
			Email:    req.Email,
			Password: hashed,
		},
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) Login(req *request.UserLogin) (*response.UserLogin, error) {
	user, err := u.userRepo.GetByEmail(util.DB, req.Email)
	if err != nil {
		return nil, err
	}

	if checkPassword(user.Password, req.Password) {
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
