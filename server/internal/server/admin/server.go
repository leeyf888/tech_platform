package admin

import (
	"context"
	mymd5 "github.com/leeyf888/go-tools/md5"
	"tech_platform/server/internal/model/user"
	"tech_platform/server/internal/pkg/response"
	"tech_platform/server/internal/store"
	adminstore "tech_platform/server/internal/store/admin"
	userstore "tech_platform/server/internal/store/user"
	"tech_platform/server/pkg/jwtutil"
)

type Handler struct {
	jwtHelper jwtutil.JWTHelper
}

func NewHandler(helper jwtutil.JWTHelper) *Handler {
	return &Handler{jwtHelper: helper}
}

func (h Handler) AdminLogin(c context.Context, req user.LoginRequest) response.ServerResponse {
	store1 := store.FromContext(c)
	u := user.User{
		Username: req.Username,
		Password: mymd5.Encryption(req.Password),
	}
	us, err := userstore.Login(store1, u)
	if err != nil {
		return response.CreateByErrorCodeMessage(response.LoginErrCode)
	}
	if us.Status == user.Forbidden {
		return response.CreateByErrorCodeMessage(response.StatusForbiddenCode)
	}

	userId, err := adminstore.AdminLogin(store1, us.UserId)
	if err != nil || userId == "" {
		return response.CreateByErrorCodeMessage(response.AdminErrCode)
	}
	token, err := h.jwtHelper.GenAdminToken(userId,true)
	if err != nil {
		return response.CreateByErrorCodeMessage(response.TokenGenErrCode)
	}

	ur := user.LoginResponse{}
	ur.UserId = userId
	ur.Status = us.Status
	ur.Token = token

	return response.CreateBySuccessData(ur)
}