package logic

import (
	"context"
	logx_cus "zerorequest/pkg"
	"zerorequest/rpc/user/internal/svc"
	___pb "zerorequest/rpc/user/proto/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserLogic) AddUser(in *___pb.UserMsg) (*___pb.Response, error) {
	// todo: add your logic here and delete this line
	//fmt.Printf("AddUser: %v\n", in)
	logger := logx_cus.GetLogger().WithContext(l.ctx)
	logger.Info("执行AddUser。。。", in)

	return &___pb.Response{
		Ok: true,
	}, nil
}
