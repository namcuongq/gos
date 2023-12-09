package ssh_bruteforce

import (
	"flag"
	"os"
	"time"
)

var (
	hostFile   string
	userFile   string
	passFile   string
	host       string
	user       string
	password   string
	port       int
	concurrent int
	sock5      string
	output     string
	cmd        string
	debug      bool
	timer      time.Duration
	flagCmd    *flag.FlagSet
)

func FlagsSetup() {
	flagCmd = flag.NewFlagSet("ssh_bruteforce", flag.ExitOnError)
	flagCmd.StringVar(&hostFile, "H", "", "File containing target hostnames or IP addresses")
	flagCmd.StringVar(&userFile, "U", "", "File containing usernames to brute force")
	flagCmd.StringVar(&passFile, "P", "", "File containing passwords to brute force")
	flagCmd.StringVar(&host, "h", "", "Target hostname or IP address")
	flagCmd.StringVar(&user, "u", "", "User to brute force")
	flagCmd.StringVar(&password, "p", "", "Password to brute force")
	flagCmd.IntVar(&port, "port", 22, "Port to brute force")
	flagCmd.IntVar(&concurrent, "c", 10, "Concurrency/threads level")
	flagCmd.StringVar(&sock5, "sock5", "", "Sock5 proxy address")
	flagCmd.StringVar(&output, "o", "success.txt", "Output file")
	flagCmd.StringVar(&cmd, "x", "", "execute command after ssh")
	flagCmd.BoolVar(&debug, "debug", false, "debug mode")
	flagCmd.DurationVar(&timer, "timer", 300*time.Millisecond, "Set timeout to ssh dial response")
	flagCmd.Parse(os.Args[2:])
}
