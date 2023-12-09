package racecondition

import (
	"flag"
	"os"
)

var (
	requestFile  string
	requestFile1 string
	numWorker    int
	host         string
	https        bool
	proxy        string
	flagCmd      *flag.FlagSet
)

func FlagsSetup() {
	flagCmd = flag.NewFlagSet("race", flag.ExitOnError)
	flagCmd.StringVar(&requestFile, "r", "", "Path of request file")
	flagCmd.StringVar(&requestFile1, "r1", "", "Path of other request file")
	flagCmd.IntVar(&numWorker, "w", 5, "Number of worker. Default: 5")
	flagCmd.StringVar(&host, "host", "", "special domain, ip")
	flagCmd.BoolVar(&https, "ssl", false, "Enable https protocol")
	flagCmd.StringVar(&proxy, "p", "", "Use a http proxy to connect to the target URL")
	flagCmd.Parse(os.Args[2:])
}
