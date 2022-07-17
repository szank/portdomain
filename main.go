package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/szank/portdomain/pkg/persistence"
	"github.com/szank/portdomain/pkg/service"
)

func main() {
	fmt.Println(os.Args)
	// TODO: pass in the termination context where necessary
	// Add a mechanism to wait for things to shut down after the termination context is closed
	// This is simple(ish) but would take a some time.
	terminateCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db := persistence.NewInMemoryDatabase()

	service := service.New(db)

	serviceResult := make(chan error)
	go func() {
		err := service.Start()
		serviceResult <- err
	}()

	select {
	case <-terminateCtx.Done():
		fmt.Println("Received a termination signal, finishing")
	case err := <-serviceResult:
		if err != nil {
			fmt.Println("Error running the service: %v", err)
			os.Exit(1)
		}
	}

	// TODO: add logrus
	fmt.Println("Done")
	os.Exit(0)
}
