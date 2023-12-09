package main

import (
	"fmt"
	"os"
	"runtime"
	fakeemail "seckit/fake_email"
	ldapbruteforce "seckit/ldap_bruteforce"
	racecondition "seckit/race_condition"
	sshbruteforce "seckit/ssh_bruteforce"
)

const (
	NAME    = "gos"
	VERSION = "1.0.0"
)

func printBanner() {
	goContainerArt := []string{
		" =======     ===     =======",
		" =          =   =   = ",
		" =         =     =  =",
		" =  =====  =     =   ======",
		" =      =  =     =         =",
		" =     =    =   =          =",
		"  =======    ===    =======",
	}

	for _, line := range goContainerArt {
		fmt.Println(line)
	}

	fmt.Println("\nVersion:", VERSION)

	fmt.Println("\nUsage: " + NAME + " <option> <agrs>")
	fmt.Println("\nOptions:")
	fmt.Println("\trace\tTesting for Race Condition")
	fmt.Println("\tmail\tTesting for Email Phishing")
	fmt.Println("\tssh\tBrute-Force SSH")
	fmt.Println("\tldap\tBrute-Force LDAP")
	fmt.Println("\ttcp\tTCP Forward Port")
	fmt.Println("EXAMPLES:")
	fmt.Println("\t" + NAME + " race -r req.txt -p http://burp:8080 -w 10 -ssl")
	fmt.Println("\t" + NAME + " mail --from support@gmail.com --to victim@company.org --attach passwd.txt --attach-name \"../etc/passwd\"")
	fmt.Println("\t" + NAME + " ssh -H hosts.txt -P pass.txt -u admin")
	fmt.Println("\t" + NAME + " ldap -s 10.10.10.1 -P pass.txt -u admin -d example.com")
	fmt.Println("\t" + NAME + " tcp -s 192.168.0.10:3389 -d 10.10.10.1:3389")
	fmt.Println("\nUse \"" + NAME + " <option> --h\" for more information about a option.")
}

func main() {
	if len(os.Args) < 2 {
		printBanner()
		return
	}

	switch os.Args[1] {
	case "race":
		racecondition.FlagsSetup()
		racecondition.Start()
	case "mail":
		fakeemail.FlagsSetup()
		fakeemail.Start()
	case "ssh":
		sshbruteforce.FlagsSetup()
		sshbruteforce.Start()
	case "ldap":
		ldapbruteforce.FlagsSetup()
		ldapbruteforce.Start()
	// case "webhook":
	// 	webhook.FlagsSetup()
	// 	webhook.Start()
	default:
		printBanner()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
