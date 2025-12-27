// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

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
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Data: types.Focus{
			Id:    "1",
			Title: "标题1",
			Image: "https://img.alicdn.com/imgextra/i2/",
			Link:  "https://www.baidu.com",
		},
	}, nil
}
