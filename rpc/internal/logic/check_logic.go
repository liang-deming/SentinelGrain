package logic

import (
	"context"

	"SentinelGrain/rpc/internal/svc"
	"SentinelGrain/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckLogic) Check(in *pb.CheckRequest) (*pb.CheckResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CheckResponse{}, nil
}
