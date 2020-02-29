package main

import (
	"os"

	"github.com/regalias/atlas-api/apiserver"
)

func main() {
	os.Exit(apiserver.Run(os.Args[1:]))
}
