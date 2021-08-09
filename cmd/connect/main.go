package main

import (
	"context"
	"im/config"
	"im/internal/connect"
	"im/pkg/db"
	"im/pkg/interceptor"
	"im/pkg/logger"
	"im/pkg/pb"
	"im/pkg/rpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化配置
	config.Init("config.yaml")

	// 初始化日志
	logger.Init(config.GetConnectServer().LogFilePath, config.GetConnectServer().LogTarget, config.GetConnectServer().LogTarget)

	db.InitRedis(config.GetRedis())

	// 初始化rpc client
	rpc.InitLogicIntClient(config.GetRpcAddr().LogicServerAddr)

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(config.GetConnectServer().LocalWsAddr)
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("connect_interceptor", nil)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 0)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop start", zap.Any("signal", s))
		_, _ = rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.GetConnectServer().LocalAddr})
		logger.Logger.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServer(server, &connect.ConnIntServer{})

	listener, err := net.Listen("tcp", config.GetConnectServer().LocalAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("Serve", zap.Error(err))
	}
}
