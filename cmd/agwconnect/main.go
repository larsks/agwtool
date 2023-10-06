package main

import (
	"agwtool/internal"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/la5nta/wl2k-go/transport/ax25/agwpe"
	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
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
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] target_callsign\n", os.Args[0])
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
		log.Fatal("missing target callsign")
	}

	targetCallsign := flag.Arg(0)

	var oldTermattr, newTermattr unix.Termios
	if err := termios.Tcgetattr(os.Stdout.Fd(), &oldTermattr); err != nil {
		log.Fatal("failed to get terminal attributes")
	}
	defer func() {
		termios.Tcsetattr(os.Stdout.Fd(), termios.TCSANOW, &oldTermattr)
	}()

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

	log.Printf("establishing connecting with %s", targetCallsign)
	conn, err := port.DialContext(context.TODO(), targetCallsign)
	if err != nil {
		panic(err)
	}

	newTermattr = oldTermattr
	termios.Cfmakeraw(&newTermattr)
	if err := termios.Tcsetattr(os.Stdout.Fd(), termios.TCSANOW, &newTermattr); err != nil {
		log.Fatalf("failed to set terminal attributes")
	}

	go func() { io.Copy(conn, os.Stdout) }()
	io.Copy(os.Stdout, conn)
}
