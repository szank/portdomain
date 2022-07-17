package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/szank/portdomain/pkg/loader"
	"github.com/szank/portdomain/pkg/persistence"
	"github.com/szank/portdomain/pkg/service"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("the service expects only one argument, the input file name")
	}

	// TODO: pass in the termination context where necessary
	// Add a mechanism to wait for things to shut down after the termination context is closed
	// This is simple(ish) but would take a some time.
	terminateCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// cannot defer file.Close if using os.Exit to spcecify the (error) exit code.
	// This can be done in a cleaner way, not sure if I got time to write it.
	file, err := loader.NewFile(os.Args[1])
	if err != nil {
		fmt.Printf("Could not load file %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	db := persistence.NewInMemoryDatabase()

	service := service.New(db, file)

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
			fmt.Printf("Error running the service: %v\n", err)
			file.Close()
			os.Exit(1)
		}
	}

	// TODO: add logrus
	file.Close()
	fmt.Println("Done")
	os.Exit(0)
}
