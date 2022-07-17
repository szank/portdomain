package main

import (
	"fmt"
	"os"

	"github.com/szank/portdomain/pkg/persistence"
	"github.com/szank/portdomain/pkg/service"
)

func main() {
	fmt.Println(os.Args)

	db := persistence.NewInMemoryDatabase()

	service := service.New(db)

	err := service.Start()
	if err != nil {
		fmt.Printf("Error running the service: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
