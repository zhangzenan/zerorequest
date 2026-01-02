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

type GetFocusByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusByIdLogic {
	return &GetFocusByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusByIdLogic) GetFocusById(req *types.FocusRequest) (resp *types.CommonResponse, err error) {
	logx.Info("GetFocusById", req.Id)
	id := req.Id
	focus := gorm.Focus{}
	err = l.svcCtx.DB.Where("id = ?", id).First(&focus).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "轮播图不存在",
			Success: false,
		}, nil
	}
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Data:    focus,
	}, nil
}
