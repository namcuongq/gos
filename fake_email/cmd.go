package fakeemail

import (
	"flag"
	"os"
)

var (
	from           string
	to             string
	subject        string
	body           string
	host           string
	port           string
	attachment     string
	attachmentName string
	flagCmd        *flag.FlagSet
)

func FlagsSetup() {
	flagCmd = flag.NewFlagSet("race", flag.ExitOnError)
	flagCmd.StringVar(&from, "from", "", "mail from")
	flagCmd.StringVar(&to, "to", "", "mail to")
	flagCmd.StringVar(&subject, "subject", "Test Subject", "subject")
	flagCmd.StringVar(&body, "body", "Test Body", "body content")
	flagCmd.StringVar(&host, "host", "", "mx record")
	flagCmd.StringVar(&port, "port", "25", "mail port")
	flagCmd.StringVar(&attachment, "attach", "", "email with attachment")
	flagCmd.StringVar(&attachmentName, "attach-name", "", "attachment file name")
	flagCmd.Parse(os.Args[2:])
}
