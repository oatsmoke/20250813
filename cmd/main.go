package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/oatsmoke/20250813/internal/handler"
	"github.com/oatsmoke/20250813/internal/repository"
	"github.com/oatsmoke/20250813/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l := service.NewLoggerService(ctx)
	defer l.Close()

	r := repository.NewTaskRepository()
	s := service.NewTaskService(ctx, r, l)
	h := handler.NewHandler(s)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h.InitRouts(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
