package main

import (
	"agwtool/internal"
	"strings"

	//"bufio"
	"context"
	"fmt"

	flag "github.com/spf13/pflag"

	//"io"
	"log"
	"net"
	"os"

	"github.com/la5nta/wl2k-go/transport/ax25/agwpe"
	//"github.com/pkg/term/termios"
	//"golang.org/x/sys/unix"
	"github.com/chzyer/readline"
)

type (
	Options struct {
		internal.CommonOptions
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
		log.Fatal("callsign is unset")
	}

	if options.TncAddress == "" {
		log.Fatal("tnc adddress is unset")
	}

	if flag.NArg() < 1 {
		log.Fatal("missing target callsign")
	}

	targetCallsign := flag.Arg(0)

	options.MyCallsign = strings.ToUpper(options.MyCallsign)
	targetCallsign = strings.ToUpper(targetCallsign)

	log.Printf("connecting to tnc at %s", options.TncAddress)
	port, err := agwpe.OpenPortTCP(options.TncAddress, 0, options.MyCallsign)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	if _, err := port.Version(); err != nil {
		log.Fatalf("AGWPE TNC initialization failed: %w", err)
	}

	log.Printf("establishing connection with %s", targetCallsign)
	conn, err := port.DialContext(context.TODO(), targetCallsign)
	if err != nil {
		panic(err)
	}

	/*
	var oldTermattr, newTermattr unix.Termios
	if err := termios.Tcgetattr(os.Stdout.Fd(), &oldTermattr); err != nil {
		log.Fatal("failed to get terminal attributes")
	}
	defer func() {
		termios.Tcsetattr(os.Stdout.Fd(), termios.TCSANOW, &oldTermattr)
	}()

	newTermattr = oldTermattr
	termios.Cfmakeraw(&newTermattr)
	if err := termios.Tcsetattr(os.Stdout.Fd(), termios.TCSANOW, &newTermattr); err != nil {
		log.Fatalf("failed to set terminal attributes")
	}
	*/

	go func() { readFromRemote(conn) }()

	rl, err := readline.New("")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			panic(err)
		}
	}

	log.Printf("all done")
}

func readFromRemote(conn net.Conn) error {
	buf := make([]byte, 256)
	for {
		nb, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read from remote: %w", err)
		}
		if nb == 0 {
			break
		}

		start := 0
		for i := range(buf[:nb]) {
			if buf[i] == '\r' {
				os.Stdout.Write(buf[start:i])
				os.Stdout.Write([]byte("\r\n"))
				start = i
			}
		}

		os.Stdout.Write(buf[start:nb])
	}

	return nil
}

func readFromStdin(conn net.Conn) error {
	log.Printf("start reading from stdin")

	char := make([]byte, 1)
	buf := make([]byte, 256)
	buflen := 0

	for {
		nbr, err := os.Stdin.Read(char)
		if err != nil {
			return fmt.Errorf("failed read from stdin: %w", err)
		}
		if nbr == 0 {
			break
		}

		buf = append(buf, char[0])
		buflen += 1
		if buflen == 256 || char[0] == '\r' || char[0] == '\n' {
			log.Printf("writing to remote: %s", buf[:buflen])
			_, err := conn.Write(buf[:buflen])
			if err != nil {
				return fmt.Errorf("failed write to remote: %w", err)
			}
			buflen = 0
		}
	}

	log.Printf("closing connection")
	conn.Close()

	return nil
}
