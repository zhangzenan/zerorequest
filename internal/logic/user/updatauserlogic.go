// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"zerorequest/model/gorm"

	"zerorequest/internal/svc"
	"zerorequest/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdataUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdataUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdataUserLogic {
	return &UpdataUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdataUserLogic) UpdataUser(req *types.UpdateUserRequest) (resp *types.CommonResponse, err error) {
	user := gorm.User{
		Id:   req.Id,
		Name: req.Name,
		Age:  req.Age,
		Sex:  req.Sex,
	}
	err = l.svcCtx.DB.Model(&gorm.User{}).Where("id = ?", req.Id).Updates(&user).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "更新失败",
		}, nil
	}
	return &types.CommonResponse{
		Code:    200,
		Message: "更新成功",
	}, nil
}
