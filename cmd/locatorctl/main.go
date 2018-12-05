package main

import (
	"flag"
	"fmt"

	"github.com/v3io/locator/pkg/locator"
)

func main() {
	config := &locator.Config{}
	flag.IntVar(&config.Port, "port", 8080, "Listen port (default: 8080)")
	flag.StringVar(&config.Namespace, "namespace", "", "Namespace to monitor")
	flag.Parse()

	if err := locator.RunServer(config); err != nil {
		fmt.Println("Error while serving registry", err)
	}
}
