package main

import (
	"flag"
	"fmt"
	"keywords_filter/filter_server/interceptor"
	"keywords_filter/filter_server/server"
	"keywords_filter/pkg/config"
	"keywords_filter/pkg/filter"
	"keywords_filter/pkg/log"
	"keywords_filter/proto/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	configFile = flag.String("config", "../dev.config.yaml", "")
	dictFile   = flag.String("dict", "../keyword.txt", "")
	formatDict = flag.Bool("format", false, "")
)

func main() {
	flag.Parse()
	if *formatDict {
		filter.OverWriteDict(*dictFile)
		return
	}
	config.InitConfig(*configFile)
	cfg := config.GetConfig()
	filter.InitFilter(*dictFile)
	Mylog := log.NewLogger()
	Mylog.SetLevel(log.Info)

	Mylog.SetOutput(log.GetRotateWriter(cfg.Log.LogPath))
	Mylog.SetPrintCaller(true)
	// log.My_log.SetCaller(caller func() (file string, line int, funcName string, err error))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.IP, cfg.Server.Port))
	if err != nil {
		Mylog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryHandler))
	service := server.NewFilterService(filter.Getfilter())
	proto.RegisterFilterServer(s, service)

	h := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, h)

	if err := s.Serve(lis); err != nil {
		Mylog.Fatalf("failed to serve: %v", err)
	}

}
