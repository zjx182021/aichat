package main

import (
	"chatgpt-web-service/chat-server/data"
	"chatgpt-web-service/chat-server/server"
	"chatgpt-web-service/chat-server/vector_data"
	"chatgpt-web-service/interceptor"
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/db/mysql"
	"chatgpt-web-service/pkg/db/redis"
	"chatgpt-web-service/pkg/db/vector"
	"chatgpt-web-service/proto"

	"chatgpt-web-service/pkg/log"
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	configFile = flag.String("config", "../dev.config.yaml", "")
	// dictFile   = flag.String("dict", "../keyword.txt", "")
	// formatDict = flag.Bool("format", false, "")
)

func main() {
	flag.Parse()

	config.InitConfig(*configFile)
	cfg := config.GetConfig()
	Mylog := log.NewLogger()
	Mylog.SetLevel(log.Info)

	Mylog.SetOutput(log.GetRotateWriter(cfg.Log.LogPath))
	Mylog.SetPrintCaller(true)
	mysql.Initmysql()
	data.Inittable()
	redis.InitRedis()
	vector.InitDB(cfg)
	// log.My_log.SetCaller(caller func() (file string, line int, funcName string, err error))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.IP, cfg.Server.Port))
	if err != nil {
		Mylog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryHandler))
	// records := data.IChatRecords.(data.ChatRecords)
	chatserver := server.NewChatService(data.Newrecords(), cfg, Mylog, vector_data.NewChatRecordsData(cfg, vector.GetVdb()))
	proto.RegisterChatServer(s, chatserver)
	h := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, h)

	if err := s.Serve(lis); err != nil {
		Mylog.Fatalf("failed to serve: %v", err)
	}

}
