package internal

import (
	"agwtool/pkg/env"
	"flag"
)

type (
	CommonOptions struct {
		MyCallsign string
		TncAddress string
	}
)

func InitCommonOptions(options *CommonOptions) {
	flag.StringVar(&options.MyCallsign, "callsign", env.Getenv("AGW_CALLSIGN", ""), "Your callsign")
	flag.StringVar(&options.TncAddress, "tncaddress", env.Getenv("AGW_TNCADDRESS", "127.0.0.1:8000"), "AGW TNC Address")
}
