package rpc

import (
	"context"
	"fmt"
	"im/pkg/grpclib"
	"im/pkg/logger"
	"im/pkg/pb"

	"google.golang.org/grpc"
)

var (
	LogicIntClient    pb.LogicIntClient
	ConnectIntClient  pb.ConnectIntClient
	BusinessIntClient pb.BusinessIntClient
)

func InitLogicIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	LogicIntClient = pb.NewLogicIntClient(conn)
}

func InitConnectIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, grpclib.Name)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	ConnectIntClient = pb.NewConnectIntClient(conn)
}
