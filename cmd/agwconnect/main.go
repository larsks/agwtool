package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/la5nta/wl2k-go/transport/ax25/agwpe"
	flag "github.com/spf13/pflag"

	"agwtool/internal"
)

type (
	Options struct {
		internal.CommonOptions

		RecvMapCrLf bool
		SendMapLfCr bool
	}
)

var options Options

func init() {
	internal.InitCommonOptions(&options.CommonOptions)

	flag.BoolVarP(&options.RecvMapCrLf, "recv-crlf", "", false, "Map CR to LF on receive")
	flag.BoolVarP(&options.SendMapLfCr, "send-lfcr", "", false, "Map LF to CR on send")
}

func printUsage() {
	fmt.Printf("Usage: %s [options] target_callsign\n", os.Args[0])
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if options.MyCallsign == "" {
		log.Fatal("local callsign is unset")
	}

	if options.TncAddress == "" {
		log.Fatal("tnc adddress is unset")
	}

	if flag.NArg() < 1 {
		log.Fatal("missing target callsign")
	}

	targetCallsign := flag.Arg(0)
	if targetCallsign == "" {
		log.Fatal("target callsign is unset")
	}

	options.MyCallsign = strings.ToUpper(options.MyCallsign)
	targetCallsign = strings.ToUpper(targetCallsign)

	log.Printf("connecting to tnc at %s", options.TncAddress)
	port, err := agwpe.OpenPortTCP(options.TncAddress, 0, options.MyCallsign)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	if version, err := port.Version(); err != nil {
		log.Fatalf("AGWPE TNC initialization failed: %v", err)
	} else {
		log.Printf("TNC version = %s", version)
	}

	log.Printf("establishing connection with %s", targetCallsign)
	conn, err := port.DialContext(context.TODO(), targetCallsign)
	if err != nil {
		panic(err)
	}

	// copy bytes from remote to stdout
	go func() { readFromRemote(conn) }()

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		if err != nil {
			panic(err)
		}

		if options.SendMapLfCr {
			_, err = conn.Write(append([]byte(line), '\r'))
		} else {
			_, err = conn.Write(append([]byte(line), '\n'))
		}

		if err != nil {
			panic(err)
		}
	}

	log.Printf("all done")
}

func readFromRemote(conn net.Conn) error {
	buf := make([]byte, 8192)
	for {
		nb, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read from remote: %w", err)
		}
		if nb == 0 {
			break
		}

		for i := range buf[:nb] {
			if buf[i] == '\r' && options.RecvMapCrLf {
				buf[i] = '\n'
			}
		}

		os.Stdout.Write(buf[:nb])
	}

	log.Printf("remote closed connection")

	return nil
}
