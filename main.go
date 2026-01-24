package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/ServerMode"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Controller"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Starter"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Router"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

func main() {
	Util, err := Util.NewUtil(State.Leader, ServerMode.Test)

	quit := make(chan os.Signal, 1)
	signalGRPCService := make(chan int, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	if err != nil {
		log.Fatal(err.Error())
	}

	mdl, err := Model.NewModel(Util)

	if err != nil {
		log.Fatal(err.Error())
	}

	ctrl, err := Controller.NewController(Util, &mdl)

	router := Router.NewRouter(ctrl)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	go GRPC_Starter.Run_GRPC_Server(Util, signalGRPCService)

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
