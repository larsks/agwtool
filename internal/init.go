package internal

import (
	"agwtool/pkg/env"

	flag "github.com/spf13/pflag"
)

type (
	CommonOptions struct {
		WorkingDirectory string
		MyCallsign       string
		TncAddress       string
	}
)

func InitCommonOptions(options *CommonOptions) {
	flag.StringVarP(&options.MyCallsign, "callsign", "c", env.Getenv("AGW_CALLSIGN", ""), "Your callsign")
	flag.StringVarP(&options.TncAddress, "tncaddress", "t", env.Getenv("AGW_TNCADDRESS", "127.0.0.1:8000"), "AGW TNC Address")
}
