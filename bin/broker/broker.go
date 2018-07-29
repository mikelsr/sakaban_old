package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/mikelsr/sakaban/broker"
)

const maxTCP = 65535

func main() {
	addr := flag.String("addr", broker.HttpDefaultAddr, "Listening address, ommit for default")
	port := flag.Int("port", broker.HttpDefaultPort, "Listening port, ommit for default")
	flag.Parse()

	if *port < 1 || *port > maxTCP {
		fmt.Printf("Error: Invalid TCP port '%d'\n", *port)
		flag.PrintDefaults()
		os.Exit(1)
	}

	b := broker.NewBroker()
	err := b.ListenAndServe(*addr, *port)
	if err != nil {
		panic(err)
	}
}
