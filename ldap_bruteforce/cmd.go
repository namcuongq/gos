package ldap_bruteforce

import (
	"flag"
	"os"
)

var (
	ldapServer string
	ldapUser   string
	ldapPass   string
	ldapDomain string

	ldapUserFile string
	ldapPassFile string

	outFile string
	flagCmd *flag.FlagSet
)

func FlagsSetup() {
	flagCmd = flag.NewFlagSet("ssh_bruteforce", flag.ExitOnError)
	flagCmd.StringVar(&ldapDomain, "d", "", "Domain")
	flagCmd.StringVar(&ldapUser, "u", "", "Single username")
	flagCmd.StringVar(&ldapPass, "p", "", "Single password")
	flagCmd.StringVar(&ldapServer, "s", "", "LDAP Server")
	flagCmd.StringVar(&ldapUserFile, "U", "", "Users.txt file")
	flagCmd.StringVar(&ldapPassFile, "P", "", "Password.txt file")
	flagCmd.StringVar(&outFile, "f", "success.txt", "Output file")
	flagCmd.Parse(os.Args[2:])
}
