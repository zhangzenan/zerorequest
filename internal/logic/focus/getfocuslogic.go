// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package focus

import (
	"context"

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
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data: []types.Focus{
			{
				Id:    "1",
				Title: "标题1",
				Image: "https://img.alicdn.com/imgextra/i2/",
				Link:  "https://www.baidu.com",
			},
			{
				Id:    "2",
				Title: "标题2",
				Image: "https://img.alicdn.com/imgextra/i2/",
				Link:  "https://www.baidu.com",
			},
		},
	}, nil
}
