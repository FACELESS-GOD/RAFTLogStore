package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/ServerMode"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Controller"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Starter"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Router"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/ServiceRegistry"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

func main() {
	mu := sync.Mutex{}
	Util, err := Util.NewUtil(State.Follower, ServerMode.Test, &mu)

	quit := make(chan os.Signal, 1)
	defer close(quit)

	signalGRPCService := make(chan int, 1)
	defer close(signalGRPCService)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	if err != nil {
		log.Fatal(err.Error())
	}

	logchan := make(chan Log.LogStuct)
	defer close(logchan)

	res := make(chan bool, 1)
	defer close(res)

	mdl, err := Model.NewModel(Util, &mu)
	if err != nil {
		log.Fatal(err.Error())
	}
	mdl.AddLogChan = logchan	

	grpcService, err := GRPC_Starter.NewGRPCService(logchan, res, Util, mdl, &mu)

	ctrl, err := Controller.NewController(Util, &mdl)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := Router.NewRouter(ctrl)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	service, err := ServiceRegistry.RegisterService(Util, mdl)
	if err != nil {
		log.Fatal(err.Error())
	}
	grpcService.Service = service
	
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	go grpcService.RecurAddLog()

	go grpcService.ElectionsCheck()

	go grpcService.Run_GRPC_Server(Util, signalGRPCService)

	<-quit

	log.Println("Shutdown signal received: server shutting down...")

	signalGRPCService <- 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
