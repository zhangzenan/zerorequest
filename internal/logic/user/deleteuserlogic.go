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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UserRequest) (resp *types.CommonResponse, err error) {
	id := req.Id
	user := gorm.User{
		Id: id,
	}
	err = l.svcCtx.DB.Delete(&user).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    500,
			Data:    nil,
			Message: "删除失败",
			Success: false,
		}, nil
	}
	return &types.CommonResponse{
		Code:    200,
		Data:    nil,
		Message: "删除成功",
		Success: true,
	}, nil
}
