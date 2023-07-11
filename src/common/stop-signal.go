package common

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignal() {
	stop := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sig)
		<-sig
		fmt.Println("got interrupt, shutting down...")
		stop <- struct{}{}
	}()
	<-stop
}
