package main

import (
	"im/config"
	"im/internal/logic/api"
	"im/pkg/db"
	"im/pkg/interceptor"
	"im/pkg/logger"
	"im/pkg/pb"
	"im/pkg/rpc"
	"im/pkg/urlwhitelist"
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
	logger.Init(config.GetLogicServer().LogFilePath, config.GetLogicServer().LogTarget, config.GetLogicServer().LogTarget)

	// 初始化数据库
	db.InitMysql(config.GetMysql())
	db.InitRedis(config.GetRedis())

	// 初始化内部rpc client
	rpc.InitConnectIntClient(config.GetRpcAddr().ConnectServerAddr)

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("logic_int_interceptor", urlwhitelist.Logic)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 0)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterLogicIntServer(server, &api.LogicIntServer{})
	pb.RegisterLogicExtServer(server, &api.LogicExtServer{})

	listen, err := net.Listen("tcp", config.GetLogicServer().LocalAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("Serve error", zap.Error(err))
	}
}
