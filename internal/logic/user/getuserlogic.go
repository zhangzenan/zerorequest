// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"zerorequest/internal/svc"
	"zerorequest/internal/types"
	"zerorequest/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser() (resp *types.CommonResponse, err error) {
	logx.Info("执行GetUser。。。")
	user := []gorm.User{}
	//l.svcCtx.DB.Debug().Find(&user)

	//执行原生sql
	l.svcCtx.DB.Raw("select * from user").Scan(&user)
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Data:    user,
	}, nil
}
