package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	group, cancelCtx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		return start(cancelCtx, ":8000", &httpHandler{})
	})
	group.Go(func() error {
		return stopMonitor(cancelCtx)
	})
	if err := group.Wait(); err != nil {
		fmt.Errorf("group return err: %+v", err)
	}

	fmt.Println("Shutdown.")
}

func stopMonitor(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal routine：other work done")
		case msg := <-ch:
			fmt.Printf("get a signal: %s\n", msg.String())
			time.Sleep(5 * time.Second)
			return fmt.Errorf("quit")
		}
	}
}

func start(ctx context.Context, addr string, h http.Handler) error {
	server := http.Server{
		Addr:    addr,
		Handler: h,
	}

	go func(ctx context.Context) {
		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		<-ctx.Done()
		fmt.Printf("http server %s ctx done\n", server.Addr)
		if err := server.Shutdown(ctx1); err != nil {
			fmt.Printf("http server %s shutdown err : %s\n", server.Addr, err)
		}
	}(ctx)
	fmt.Println("http server start!")
	return server.ListenAndServe()
}

type httpHandler struct {
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}
