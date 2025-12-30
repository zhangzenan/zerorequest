// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package focus

import (
	"context"
	"zerorequest/model/gorm"

	"zerorequest/internal/svc"
	"zerorequest/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFocusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusLogic {
	return &GetFocusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusLogic) GetFocus() (resp *types.CommonResponse, err error) {
	//获取所有数据,定义一个切片
	focus := []gorm.Focus{}
	l.svcCtx.DB.Find(&focus)
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data:    focus,
	}, nil
}
