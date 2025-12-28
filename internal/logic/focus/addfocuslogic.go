// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package focus

import (
	"context"

	"zerorequest/internal/svc"
	"zerorequest/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFocusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFocusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFocusLogic {
	return &AddFocusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFocusLogic) AddFocus(req *types.AddFocusRequest) (resp *types.CommonResponse, err error) {
	logx.Info("AddFocus", req.Title, req.Image, req.Link)
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data: types.Focus{
			Id:    "1",
			Title: req.Title,
			Image: req.Image,
			Link:  req.Link,
		},
	}, nil
}
