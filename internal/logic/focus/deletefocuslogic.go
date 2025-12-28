// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package focus

import (
	"context"

	"zerorequest/internal/svc"
	"zerorequest/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFocusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFocusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFocusLogic {
	return &DeleteFocusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFocusLogic) DeleteFocus(req *types.FocusRequest) (resp *types.CommonResponse, err error) {
	logx.Info("DeleteFocus", req.Id)
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
	}, nil
}
