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

	service.Start()
}
