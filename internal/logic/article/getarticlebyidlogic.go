// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"zerorequest/internal/svc"
	"zerorequest/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleByIdLogic {
	return &GetArticleByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleByIdLogic) GetArticleById(req *types.ArticleRequest) (resp *types.CommonResponse, err error) {
	//获取动态路由传值
	logx.Info("GetArticleById", req.Id)
	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data: types.Article{
			Id:      1,
			Title:   "标题",
			Content: "内容",
		},
	}, nil
}
