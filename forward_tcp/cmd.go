package forwardtcp

import (
	"flag"
	"os"
)

var (
	src string
	dst string

	flagCmd *flag.FlagSet
)

func FlagsSetup() {
	flagCmd = flag.NewFlagSet("forward_tcp", flag.ExitOnError)
	flagCmd.StringVar(&src, "s", "", "source")
	flagCmd.StringVar(&dst, "d", "", "destination")
	flagCmd.Parse(os.Args[2:])
}
