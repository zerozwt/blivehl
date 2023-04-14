package handler

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/service"
	"github.com/zerozwt/blivehl/server/utils"
)

var notLoginRsp []byte = []byte(`{"code":114514,"msg":"not logged in","data":{}}`)

func init() {
	engine.RegisterApi("/user/login", login)
	engine.RegisterApi("/user/logout", logout, loginChecker)
	engine.RegisterApi("/user/info", userInfo, loginChecker)
	engine.RegisterApi("/user/pass", changePassword, loginChecker)
}

func loginChecker(ctx *engine.Context) {
	writeNotLogin := func() {
		ctx.RawResponse.Header().Set("Content-Type", "application/json")
		ctx.RawResponse.Write(notLoginRsp)
	}

	cookie, ok := utils.GetCtxCookie(ctx)
	if !ok {
		writeNotLogin()
		return
	}

	user, isAdmin, err := service.GetLoginService().CheckCookie(cookie)
	if err != nil {
		writeNotLogin()
		return
	}

	utils.PutCtxUser(ctx, user)
	if isAdmin {
		utils.SetCtxAdmin(ctx, isAdmin)
	}

	ctx.Next()
}

func login(ctx *engine.Context, req *bs.LoginRequest) (*bs.LoginResponse, error) {
	cookie, err := service.GetLoginService().LoginAndGenerateCookie(req.User, req.Password)
	if err != nil {
		return nil, err
	}

	utils.PutCtxCookie(ctx, cookie)
	return &bs.LoginResponse{}, nil
}

func logout(ctx *engine.Context, req *bs.LoginResponse) (*bs.LoginResponse, error) {
	cookie, _ := utils.GetCtxCookie(ctx)
	err := service.GetLoginService().DeleteCookie(cookie)
	utils.PutCtxCookie(ctx, "-")
	return &bs.LoginResponse{}, err
}

func userInfo(ctx *engine.Context, req *bs.LoginResponse) (*bs.UserInfoResponse, error) {
	userID, _ := utils.GetCtxUser(ctx)

	userName, isAdmin, err := service.GetLoginService().GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	return &bs.UserInfoResponse{Name: userName, IsAdmin: isAdmin}, nil
}

func changePassword(ctx *engine.Context, req *bs.ChangePasswordRequest) (*bs.LoginResponse, error) {
	userID, _ := utils.GetCtxUser(ctx)
	return &bs.LoginResponse{}, service.GetLoginService().ChangePassword(userID, req.OldPass, req.NewPass)
}
