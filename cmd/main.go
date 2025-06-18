package main

import (
	"context"
	"fmt"
	xloader "github.com/mhthrh/common_pkg/pkg/loader"
	l "github.com/mhthrh/common_pkg/pkg/logger"
	cnfg "github.com/mhthrh/common_pkg/pkg/model/config"
	user "github.com/mhthrh/common_pkg/pkg/model/user/grpc/v1"
	"github.com/mhthrh/common_pkg/pkg/xErrors"
	"github.com/mhthrh/common_pkg/util/generic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"restfullApi/pkg/gRPC"
	"syscall"
)

const (
	configPath = "src/usersrv/config/file"
	appName    = "user"
	url        = "https://vault.mhthrh.co.ca"
	secret     = "AnKoloft@~delNazok!12345"
	logName    = "x-app.user.service"
)

var (
	osInterrupt       chan os.Signal
	internalInterrupt chan *xErrors.Error
)

func init() {
	osInterrupt = make(chan os.Signal)
	internalInterrupt = make(chan *xErrors.Error)
	signal.Notify(osInterrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := l.NewLogger(logName)
	defer func(){
		cancel()
		logger.LogSync()
	}() 

	logger.Info(ctx, "Loading config...")

	config, err := xloader.New(url, configPath, "", "", secret, true)
	if err != nil {
		logger.Fatal(ctx, "config loader error", zap.Any("config loader failed", err))
	}
	err = config.Read()
	if err != nil {
		logger.Fatal(ctx, "config reader error", zap.Any("config loader failed", err))
	}
	logger.Info(ctx, "user service configs loaded successfully")
	grpcs, err := config.GetGrpcs()
	if err != nil {
		logger.Fatal(ctx, "config loader error", zap.Any("config loader failed", err))
	}
	g := generic.Filter(grpcs, appName, func(t cnfg.Grpc, s string) bool {
		if t.Srv == appName {
			return true
		}
		return false
	})

	lis, e := net.Listen("tcp", fmt.Sprintf("%s:%d", g.Host, g.Port))
	if e != nil {
		log.Fatalf("Error listening on %s. error %v", g.Host, e)
	}
	rpcServer := grpc.NewServer()
	defer rpcServer.GracefulStop()

	user.RegisterUserServiceServer(
		rpcServer, gRPC.New(logger),
	)

	go func() {
		defer log.Println("listener has been stopped")

		log.Println("starting listener...")
		reflection.Register(rpcServer)
		if e = rpcServer.Serve(lis); e != nil {
			internalInterrupt <- err
		}
	}()

	logger.Info(ctx, "gRPC server started on port ", zap.Int("grpc port", g.Port))

	logger.Info(ctx, "service listening for any interrupt signals...")

	select {
	case <-osInterrupt:
		logger.Info(ctx, "OS interrupt signal received")
		rpcServer.Stop()
	case e := <-internalInterrupt:
		logger.Error(ctx, "user service listener interrupt signal received, %+v", zap.Any("lis", e))
	}

	logger.Info(ctx, "stopping user service...")

}
