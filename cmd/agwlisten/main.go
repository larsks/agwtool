package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"agwtool/internal"

	"github.com/creack/pty"
	"github.com/la5nta/wl2k-go/transport/ax25/agwpe"
)

type (
	Options struct {
		internal.CommonOptions
	}
)

var options Options

func init() {
	internal.InitCommonOptions(&options.CommonOptions)
}

func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] command [arg [...]]\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if options.MyCallsign == "" {
		log.Fatal("callsign is unset")
	}

	if options.TncAddress == "" {
		log.Fatal("tnc adddress is unset")
	}

	if flag.NArg() < 1 {
		log.Fatal("missing command")
	}

	log.Printf("connecting to tnc at %s", options.TncAddress)
	tnc, err := agwpe.OpenTCP(options.TncAddress)
	if err != nil {
		panic(err)
	}
	defer tnc.Close()

	port, err := tnc.RegisterPort(0, options.MyCallsign)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	listener, err := port.Listen()
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("waiting for connection")
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("connection from %s", conn.RemoteAddr())

		log.Printf("starting %s", flag.Arg(0))
		cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
		fd, err := pty.Start(cmd)
		go func() { io.Copy(fd, conn) }()
		io.Copy(conn, fd)

		conn.Close()
	}
}
